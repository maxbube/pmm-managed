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

package management

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/api/managementpb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
)

// MongoDBService MongoDB Management Service.
//nolint:unused
type MongoDBService struct {
	db       *reform.DB
	registry agentsRegistry
}

// NewMongoDBService creates new MySQL Management Service.
func NewMongoDBService(db *reform.DB, registry agentsRegistry) *MongoDBService {
	return &MongoDBService{db, registry}
}

// Add adds "MongoDB Service", "MongoDB Exporter Agent" and "QAN MongoDB Profiler".
func (s *MongoDBService) Add(ctx context.Context, req *managementpb.AddMongoDBRequest) (*managementpb.AddMongoDBResponse, error) {
	res := new(managementpb.AddMongoDBResponse)

	if e := s.db.InTransaction(func(tx *reform.TX) error {
		service, err := models.AddNewService(tx.Querier, models.MongoDBServiceType, &models.AddDBMSServiceParams{
			ServiceName:    req.ServiceName,
			NodeID:         req.NodeId,
			Environment:    req.Environment,
			Cluster:        req.Cluster,
			ReplicationSet: req.ReplicationSet,
			Address:        pointer.ToStringOrNil(req.Address),
			Port:           pointer.ToUint16OrNil(uint16(req.Port)),
			CustomLabels:   req.CustomLabels,
		})
		if err != nil {
			return err
		}

		invService, err := services.ToAPIService(service)
		if err != nil {
			return err
		}
		res.Service = invService.(*inventorypb.MongoDBService)

		row, err := models.CreateAgent(tx.Querier, models.MongoDBExporterType, &models.CreateAgentParams{
			PMMAgentID: req.PmmAgentId,
			ServiceID:  service.ServiceID,
			Username:   req.Username,
			Password:   req.Password,
		})
		if err != nil {
			return err
		}

		if !req.SkipConnectionCheck {
			if err = s.registry.CheckConnectionToService(ctx, service, row); err != nil {
				return err
			}
		}

		agent, err := services.ToAPIAgent(tx.Querier, row)
		if err != nil {
			return err
		}
		res.MongodbExporter = agent.(*inventorypb.MongoDBExporter)

		if req.QanMongodbProfiler {
			row, err = models.CreateAgent(tx.Querier, models.QANMongoDBProfilerAgentType, &models.CreateAgentParams{
				PMMAgentID: req.PmmAgentId,
				ServiceID:  service.ServiceID,
				Username:   req.Username,
				Password:   req.Password,
			})
			if err != nil {
				return err
			}

			agent, err := services.ToAPIAgent(tx.Querier, row)
			if err != nil {
				return err
			}
			res.QanMongodbProfiler = agent.(*inventorypb.QANMongoDBProfilerAgent)
		}

		return nil
	}); e != nil {
		return nil, e
	}

	s.registry.SendSetStateRequest(ctx, req.PmmAgentId)
	return res, nil
}