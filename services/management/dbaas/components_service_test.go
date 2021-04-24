// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package dbaas

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	goversion "github.com/hashicorp/go-version"
	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

const versionServiceURL = "https://check.percona.com/versions/v1"

func TestComponentService(t *testing.T) {
	const (
		clusterName = "pxcCluster"
		kubeConfig  = "{}"
	)

	setup := func(t *testing.T) (ctx context.Context, cs dbaasv1beta1.ComponentsServer, dbaasClient *mockDbaasClient) {
		t.Helper()

		ctx = logger.Set(context.Background(), t.Name())
		uuid.SetRand(new(tests.IDReader))

		sqlDB := testdb.Open(t, models.SetupFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		dbaasClient = new(mockDbaasClient)

		kubernetesCluster, err := models.CreateKubernetesCluster(db.Querier, &models.CreateKubernetesClusterParams{
			KubernetesClusterName: clusterName,
			KubeConfig:            kubeConfig,
		})
		require.NoError(t, err)

		t.Cleanup(func() {
			uuid.SetRand(nil)
			dbaasClient.AssertExpectations(t)
			assert.NoError(t, db.Delete(kubernetesCluster))
			require.NoError(t, sqlDB.Close())
		})

		vsc := NewVersionServiceClient(versionServiceURL)
		cs = NewComponentsService(db, dbaasClient, vsc)

		return
	}

	t.Run("PXC", func(t *testing.T) {
		t.Run("BasicGet", func(t *testing.T) {
			ctx, cs, dbaasClient := setup(t)

			dbaasClient.On("CheckKubernetesClusterConnection", mock.Anything, "{}").Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
				Operators: &controllerv1beta1.Operators{Xtradb: &controllerv1beta1.Operator{
					Status:  controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK,
					Version: "1.7.0",
				}},
				Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
			}, nil)

			pxcComponents, err := cs.GetPXCComponents(ctx, &dbaasv1beta1.GetPXCComponentsRequest{
				KubernetesClusterName: clusterName,
			})
			require.NoError(t, err)
			require.NotNil(t, pxcComponents)

			expected := map[string]*dbaasv1beta1.Component{
				"8.0.19-10.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false},
				"8.0.20-11.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false},
				"8.0.20-11.2": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false},
				"8.0.21-12.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false, Default: true},
			}
			require.Equal(t, 1, len(pxcComponents.Versions))
			assert.Equal(t, expected, pxcComponents.Versions[0].Matrix.Pxc)
		})

		t.Run("Change", func(t *testing.T) {
			ctx, cs, dbaasClient := setup(t)

			dbaasClient.On("CheckKubernetesClusterConnection", mock.Anything, "{}").Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
				Operators: &controllerv1beta1.Operators{Xtradb: &controllerv1beta1.Operator{
					Status:  controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK,
					Version: "1.7.0",
				}},
				Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
			}, nil)

			resp, err := cs.ChangePXCComponents(ctx, &dbaasv1beta1.ChangePXCComponentsRequest{
				KubernetesClusterName: clusterName,
				Pxc: &dbaasv1beta1.ChangeComponent{
					DefaultVersion: "8.0.19-10.1",
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "8.0.20-11.1",
						Disable: true,
					}, {
						Version: "8.0.20-11.2",
						Disable: true,
					}},
				},
				Proxysql: nil,
			})
			require.NoError(t, err)
			require.NotNil(t, resp)

			pxcComponents, err := cs.GetPXCComponents(ctx, &dbaasv1beta1.GetPXCComponentsRequest{
				KubernetesClusterName: clusterName,
			})
			require.NoError(t, err)
			require.NotNil(t, pxcComponents)

			expected := map[string]*dbaasv1beta1.Component{
				"8.0.19-10.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false, Default: true},
				"8.0.20-11.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false, Disabled: true},
				"8.0.20-11.2": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false, Disabled: true},
				"8.0.21-12.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false},
			}
			require.Equal(t, 1, len(pxcComponents.Versions))
			assert.Equal(t, expected, pxcComponents.Versions[0].Matrix.Pxc)

			t.Run("Change Again", func(t *testing.T) {
				resp, err := cs.ChangePXCComponents(ctx, &dbaasv1beta1.ChangePXCComponentsRequest{
					KubernetesClusterName: clusterName,
					Pxc: &dbaasv1beta1.ChangeComponent{
						DefaultVersion: "8.0.20-11.1",
						Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
							Version: "8.0.20-11.1",
							Enable:  true,
						}},
					},
					Proxysql: nil,
				})
				require.NoError(t, err)
				require.NotNil(t, resp)

				pxcComponents, err := cs.GetPXCComponents(ctx, &dbaasv1beta1.GetPXCComponentsRequest{
					KubernetesClusterName: clusterName,
				})
				require.NoError(t, err)
				require.NotNil(t, pxcComponents)

				expected := map[string]*dbaasv1beta1.Component{
					"8.0.19-10.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false},
					"8.0.20-11.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false, Default: true},
					"8.0.20-11.2": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false, Disabled: true},
					"8.0.21-12.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false},
				}
				require.Equal(t, 1, len(pxcComponents.Versions))
				assert.Equal(t, expected, pxcComponents.Versions[0].Matrix.Pxc)
			})
		})

		t.Run("Don't let disable and make default same version", func(t *testing.T) {
			ctx, cs, _ := setup(t)

			resp, err := cs.ChangePXCComponents(ctx, &dbaasv1beta1.ChangePXCComponentsRequest{
				KubernetesClusterName: clusterName,
				Pxc: &dbaasv1beta1.ChangeComponent{
					DefaultVersion: "8.0.19-10.1",
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "8.0.19-10.1",
						Disable: true,
						Enable:  false,
					}},
				},
				Proxysql: nil,
			})
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, fmt.Sprintf("default version can't be disabled, cluster: %s, component: pxc", clusterName)), err)
			require.Nil(t, resp)
		})

		t.Run("enable and disable", func(t *testing.T) {
			ctx, cs, _ := setup(t)

			resp, err := cs.ChangePXCComponents(ctx, &dbaasv1beta1.ChangePXCComponentsRequest{
				KubernetesClusterName: clusterName,
				Pxc:                   nil,
				Proxysql: &dbaasv1beta1.ChangeComponent{
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "8.0.19-10.1",
						Disable: true,
						Enable:  true,
					}},
				},
			})
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, fmt.Sprintf("enable and disable for version 8.0.19-10.1 can't be passed together, cluster: %s, component: proxySQL", clusterName)), err)
			require.Nil(t, resp)
		})
	})

	t.Run("PSMDB", func(t *testing.T) {
		t.Run("BasicGet", func(t *testing.T) {
			ctx, cs, dbaasClient := setup(t)

			dbaasClient.On("CheckKubernetesClusterConnection", mock.Anything, "{}").Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
				Operators: &controllerv1beta1.Operators{Psmdb: &controllerv1beta1.Operator{
					Status:  controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK,
					Version: "1.6.0",
				}},
				Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
			}, nil)

			psmdbComponents, err := cs.GetPSMDBComponents(ctx, &dbaasv1beta1.GetPSMDBComponentsRequest{
				KubernetesClusterName: clusterName,
			})
			require.NoError(t, err)
			require.NotNil(t, psmdbComponents)

			expected := map[string]*dbaasv1beta1.Component{
				"4.2.7-7":   {ImagePath: "percona/percona-server-mongodb:4.2.7-7", ImageHash: "1d8a0859b48a3e9cadf9ad7308ec5aa4b278a64ca32ff5d887156b1b46146b13", Status: "available", Critical: false},
				"4.2.8-8":   {ImagePath: "percona/percona-server-mongodb:4.2.8-8", ImageHash: "a66e889d3e986413e41083a9c887f33173da05a41c8bd107cf50eede4588a505", Status: "available", Critical: false},
				"4.2.11-12": {ImagePath: "percona/percona-server-mongodb:4.2.11-12", ImageHash: "1909cb7a6ecea9bf0535b54aa86b9ae74ba2fa303c55cf4a1a54262fb0edbd3c", Status: "recommended", Critical: false},
				"4.4.2-4":   {ImagePath: "percona/percona-server-mongodb:4.4.2-4", ImageHash: "991d6049059e5eb1a74981290d829a5fb4ab0554993748fde1e67b2f46f26bf0", Status: "recommended", Critical: false, Default: true},
			}
			require.Equal(t, 1, len(psmdbComponents.Versions))
			assert.Equal(t, expected, psmdbComponents.Versions[0].Matrix.Mongod)
		})

		t.Run("Change", func(t *testing.T) {
			ctx, cs, dbaasClient := setup(t)

			dbaasClient.On("CheckKubernetesClusterConnection", mock.Anything, "{}").Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
				Operators: &controllerv1beta1.Operators{Psmdb: &controllerv1beta1.Operator{
					Status:  controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK,
					Version: "1.6.0",
				}},
				Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
			}, nil)

			resp, err := cs.ChangePSMDBComponents(ctx, &dbaasv1beta1.ChangePSMDBComponentsRequest{
				KubernetesClusterName: clusterName,
				Mongod: &dbaasv1beta1.ChangeComponent{
					DefaultVersion: "4.2.8-8",
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "4.2.7-7",
						Disable: true,
					}, {
						Version: "4.4.2-4",
						Disable: true,
					}},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, resp)

			psmdbComponents, err := cs.GetPSMDBComponents(ctx, &dbaasv1beta1.GetPSMDBComponentsRequest{
				KubernetesClusterName: clusterName,
			})
			require.NoError(t, err)
			require.NotNil(t, psmdbComponents)

			expected := map[string]*dbaasv1beta1.Component{
				"4.2.7-7":   {ImagePath: "percona/percona-server-mongodb:4.2.7-7", ImageHash: "1d8a0859b48a3e9cadf9ad7308ec5aa4b278a64ca32ff5d887156b1b46146b13", Status: "available", Critical: false, Disabled: true},
				"4.2.8-8":   {ImagePath: "percona/percona-server-mongodb:4.2.8-8", ImageHash: "a66e889d3e986413e41083a9c887f33173da05a41c8bd107cf50eede4588a505", Status: "available", Critical: false, Default: true},
				"4.2.11-12": {ImagePath: "percona/percona-server-mongodb:4.2.11-12", ImageHash: "1909cb7a6ecea9bf0535b54aa86b9ae74ba2fa303c55cf4a1a54262fb0edbd3c", Status: "recommended", Critical: false},
				"4.4.2-4":   {ImagePath: "percona/percona-server-mongodb:4.4.2-4", ImageHash: "991d6049059e5eb1a74981290d829a5fb4ab0554993748fde1e67b2f46f26bf0", Status: "recommended", Critical: false, Disabled: true},
			}
			require.Equal(t, 1, len(psmdbComponents.Versions))
			assert.Equal(t, expected, psmdbComponents.Versions[0].Matrix.Mongod)

			t.Run("Change Again", func(t *testing.T) {
				resp, err := cs.ChangePSMDBComponents(ctx, &dbaasv1beta1.ChangePSMDBComponentsRequest{
					KubernetesClusterName: clusterName,
					Mongod: &dbaasv1beta1.ChangeComponent{
						DefaultVersion: "4.2.11-12",
						Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
							Version: "4.4.2-4",
							Enable:  true,
						}, {
							Version: "4.2.8-8",
							Disable: true,
						}},
					},
				})
				require.NoError(t, err)
				require.NotNil(t, resp)

				psmdbComponents, err := cs.GetPSMDBComponents(ctx, &dbaasv1beta1.GetPSMDBComponentsRequest{
					KubernetesClusterName: clusterName,
				})
				require.NoError(t, err)
				require.NotNil(t, psmdbComponents)

				expected := map[string]*dbaasv1beta1.Component{
					"4.2.7-7":   {ImagePath: "percona/percona-server-mongodb:4.2.7-7", ImageHash: "1d8a0859b48a3e9cadf9ad7308ec5aa4b278a64ca32ff5d887156b1b46146b13", Status: "available", Critical: false, Disabled: true},
					"4.2.8-8":   {ImagePath: "percona/percona-server-mongodb:4.2.8-8", ImageHash: "a66e889d3e986413e41083a9c887f33173da05a41c8bd107cf50eede4588a505", Status: "available", Critical: false, Disabled: true},
					"4.2.11-12": {ImagePath: "percona/percona-server-mongodb:4.2.11-12", ImageHash: "1909cb7a6ecea9bf0535b54aa86b9ae74ba2fa303c55cf4a1a54262fb0edbd3c", Status: "recommended", Critical: false, Default: true},
					"4.4.2-4":   {ImagePath: "percona/percona-server-mongodb:4.4.2-4", ImageHash: "991d6049059e5eb1a74981290d829a5fb4ab0554993748fde1e67b2f46f26bf0", Status: "recommended", Critical: false},
				}
				require.Equal(t, 1, len(psmdbComponents.Versions))
				assert.Equal(t, expected, psmdbComponents.Versions[0].Matrix.Mongod)
			})
		})

		t.Run("Don't let disable and make default same version", func(t *testing.T) {
			ctx, cs, _ := setup(t)

			resp, err := cs.ChangePSMDBComponents(ctx, &dbaasv1beta1.ChangePSMDBComponentsRequest{
				KubernetesClusterName: clusterName,
				Mongod: &dbaasv1beta1.ChangeComponent{
					DefaultVersion: "4.2.11-12",
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "4.2.11-12",
						Disable: true,
						Enable:  false,
					}},
				},
			})
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, fmt.Sprintf("default version can't be disabled, cluster: %s, component: mongod", clusterName)), err)
			require.Nil(t, resp)
		})

		t.Run("enable and disable", func(t *testing.T) {
			ctx, cs, _ := setup(t)

			resp, err := cs.ChangePSMDBComponents(ctx, &dbaasv1beta1.ChangePSMDBComponentsRequest{
				KubernetesClusterName: clusterName,
				Mongod: &dbaasv1beta1.ChangeComponent{
					Versions: []*dbaasv1beta1.ChangeComponent_ComponentVersion{{
						Version: "4.2.11-12",
						Disable: true,
						Enable:  true,
					}},
				},
			})
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, fmt.Sprintf("enable and disable for version 4.2.11-12 can't be passed together, cluster: %s, component: mongod", clusterName)), err)
			require.Nil(t, resp)
		})
	})
}

func TestComponentServiceMatrix(t *testing.T) {
	input := map[string]componentVersion{
		"5.7.26-31.37":   {ImagePath: "percona/percona-xtradb-cluster:5.7.26-31.37", ImageHash: "9d43d8e435e4aca5c694f726cc736667cb938158635c5f01a0e9412905f1327f", Status: "available", Critical: false},
		"5.7.27-31.39":   {ImagePath: "percona/percona-xtradb-cluster:5.7.27-31.39", ImageHash: "7d8eb4d2031c32c6e96451655f359d8e5e8e047dc95bada9a28c41c158876c26", Status: "available", Critical: false},
		"5.7.28-31.41.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.28-31.41.2", ImageHash: "fccd6525aaeedb5e436e9534e2a63aebcf743c043526dd05dba8519ebddc8b30", Status: "available", Critical: true},
		"5.7.29-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.29-31.43", ImageHash: "85fb479de073770280ae601cf3ec22dc5c8cca4c8b0dc893b09503767338e6f9", Status: "available", Critical: false},
		"5.7.30-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.30-31.43", ImageHash: "b03a060e9261b37288a2153c78f86dcfc53367c36e1bcdcae046dd2d0b0721af", Status: "available", Critical: false},
		"5.7.31-31.45":   {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45", ImageHash: "3852cef43cc0c6aa791463ba6279e59dcdac3a4fb1a5616c745c1b3c68041dc2", Status: "available", Critical: false},
		"5.7.31-31.45.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45.2", ImageHash: "0decf85c7c7afacc438f5fe355dc8320ea7ffc7018ca2cb6bda3ac0c526ae172", Status: "available", Critical: false},
		"5.7.32-31.47":   {ImagePath: "percona/percona-xtradb-cluster:5.7.32-31.47", ImageHash: "7b095019ad354c336494248d6080685022e2ed46e3b53fc103b25cd12c95952b", Status: "recommended", Critical: false},
		"8.0.19-10.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false},
		"8.0.20-11.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false},
		"8.0.20-11.2":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false},
		"8.0.21-12.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false},
	}

	t.Run("All", func(t *testing.T) {
		cs := &componentsService{}
		m := cs.matrix(input, nil, nil)

		expected := map[string]*dbaasv1beta1.Component{
			"5.7.26-31.37":   {ImagePath: "percona/percona-xtradb-cluster:5.7.26-31.37", ImageHash: "9d43d8e435e4aca5c694f726cc736667cb938158635c5f01a0e9412905f1327f", Status: "available", Critical: false},
			"5.7.27-31.39":   {ImagePath: "percona/percona-xtradb-cluster:5.7.27-31.39", ImageHash: "7d8eb4d2031c32c6e96451655f359d8e5e8e047dc95bada9a28c41c158876c26", Status: "available", Critical: false},
			"5.7.28-31.41.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.28-31.41.2", ImageHash: "fccd6525aaeedb5e436e9534e2a63aebcf743c043526dd05dba8519ebddc8b30", Status: "available", Critical: true},
			"5.7.29-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.29-31.43", ImageHash: "85fb479de073770280ae601cf3ec22dc5c8cca4c8b0dc893b09503767338e6f9", Status: "available", Critical: false},
			"5.7.30-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.30-31.43", ImageHash: "b03a060e9261b37288a2153c78f86dcfc53367c36e1bcdcae046dd2d0b0721af", Status: "available", Critical: false},
			"5.7.31-31.45":   {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45", ImageHash: "3852cef43cc0c6aa791463ba6279e59dcdac3a4fb1a5616c745c1b3c68041dc2", Status: "available", Critical: false},
			"5.7.31-31.45.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45.2", ImageHash: "0decf85c7c7afacc438f5fe355dc8320ea7ffc7018ca2cb6bda3ac0c526ae172", Status: "available", Critical: false},
			"5.7.32-31.47":   {ImagePath: "percona/percona-xtradb-cluster:5.7.32-31.47", ImageHash: "7b095019ad354c336494248d6080685022e2ed46e3b53fc103b25cd12c95952b", Status: "recommended", Critical: false},
			"8.0.19-10.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false},
			"8.0.20-11.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false},
			"8.0.20-11.2":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false},
			"8.0.21-12.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false, Default: true},
		}

		assert.Equal(t, expected, m)
	})

	t.Run("Disabled and Default Components", func(t *testing.T) {
		cs := &componentsService{}

		m := cs.matrix(input, nil, &models.Component{
			DisabledVersions: []string{"8.0.20-11.2", "8.0.20-11.1"},
			DefaultVersion:   "8.0.19-10.1",
		})

		expected := map[string]*dbaasv1beta1.Component{
			"5.7.26-31.37":   {ImagePath: "percona/percona-xtradb-cluster:5.7.26-31.37", ImageHash: "9d43d8e435e4aca5c694f726cc736667cb938158635c5f01a0e9412905f1327f", Status: "available", Critical: false},
			"5.7.27-31.39":   {ImagePath: "percona/percona-xtradb-cluster:5.7.27-31.39", ImageHash: "7d8eb4d2031c32c6e96451655f359d8e5e8e047dc95bada9a28c41c158876c26", Status: "available", Critical: false},
			"5.7.28-31.41.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.28-31.41.2", ImageHash: "fccd6525aaeedb5e436e9534e2a63aebcf743c043526dd05dba8519ebddc8b30", Status: "available", Critical: true},
			"5.7.29-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.29-31.43", ImageHash: "85fb479de073770280ae601cf3ec22dc5c8cca4c8b0dc893b09503767338e6f9", Status: "available", Critical: false},
			"5.7.30-31.43":   {ImagePath: "percona/percona-xtradb-cluster:5.7.30-31.43", ImageHash: "b03a060e9261b37288a2153c78f86dcfc53367c36e1bcdcae046dd2d0b0721af", Status: "available", Critical: false},
			"5.7.31-31.45":   {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45", ImageHash: "3852cef43cc0c6aa791463ba6279e59dcdac3a4fb1a5616c745c1b3c68041dc2", Status: "available", Critical: false},
			"5.7.31-31.45.2": {ImagePath: "percona/percona-xtradb-cluster:5.7.31-31.45.2", ImageHash: "0decf85c7c7afacc438f5fe355dc8320ea7ffc7018ca2cb6bda3ac0c526ae172", Status: "available", Critical: false},
			"5.7.32-31.47":   {ImagePath: "percona/percona-xtradb-cluster:5.7.32-31.47", ImageHash: "7b095019ad354c336494248d6080685022e2ed46e3b53fc103b25cd12c95952b", Status: "recommended", Critical: false},
			"8.0.19-10.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false, Default: true},
			"8.0.20-11.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false, Disabled: true},
			"8.0.20-11.2":    {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false, Disabled: true},
			"8.0.21-12.1":    {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false},
		}

		assert.Equal(t, expected, m)
	})

	t.Run("Skip unsupported Components", func(t *testing.T) {
		cs := &componentsService{}

		minimumSupportedVersion, err := goversion.NewVersion("8.0.0")
		require.NoError(t, err)
		m := cs.matrix(input, minimumSupportedVersion, &models.Component{
			DisabledVersions: []string{"8.0.21-12.1", "8.0.20-11.1"},
			DefaultVersion:   "8.0.20-11.2",
		})

		expected := map[string]*dbaasv1beta1.Component{
			"8.0.19-10.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1", ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a", Status: "available", Critical: false},
			"8.0.20-11.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1", ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349", Status: "available", Critical: false, Disabled: true},
			"8.0.20-11.2": {ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.2", ImageHash: "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b", Status: "available", Critical: false, Default: true},
			"8.0.21-12.1": {ImagePath: "percona/percona-xtradb-cluster:8.0.21-12.1", ImageHash: "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e", Status: "recommended", Critical: false, Disabled: true},
		}

		assert.Equal(t, expected, m)
	})

	t.Run("EmptyMatrix", func(t *testing.T) {
		cs := &componentsService{}
		m := cs.matrix(map[string]componentVersion{}, nil, nil)
		assert.Equal(t, map[string]*dbaasv1beta1.Component{}, m)
	})
}

func TestFilteringOutOfUnsupportedVersions(t *testing.T) {
	t.Parallel()
	c := &componentsService{
		l:                    logrus.WithField("component", "components_service"),
		versionServiceClient: NewVersionServiceClient(versionServiceURL),
	}

	t.Run("mongod", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		params := componentsParams{
			operator:        psmdbOperator,
			operatorVersion: "1.6.0",
		}
		versions, err := c.versions(ctx, params, nil)
		require.NoError(t, err)
		parsedSupportedVersion, err := goversion.NewVersion("4.2.0")
		require.NoError(t, err)
		for _, v := range versions {
			for version := range v.Matrix.Mongod {
				parsedVersion, err := goversion.NewVersion(version)
				require.NoError(t, err)
				assert.Truef(t, parsedVersion.GreaterThanOrEqual(parsedSupportedVersion), "%s is not greater or equal to 4.2.0", version)
			}
		}
	})

	t.Run("pxc", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		params := componentsParams{
			operator:        pxcOperator,
			operatorVersion: "1.7.0",
		}
		versions, err := c.versions(ctx, params, nil)
		require.NoError(t, err)
		parsedSupportedVersion, err := goversion.NewVersion("8.0.0")
		require.NoError(t, err)
		for _, v := range versions {
			for version := range v.Matrix.Pxc {
				parsedVersion, err := goversion.NewVersion(version)
				require.NoError(t, err)
				assert.True(t, parsedVersion.GreaterThanOrEqual(parsedSupportedVersion), "%s is not greater or equal to 8.0.0", version)
			}
		}
	})
}