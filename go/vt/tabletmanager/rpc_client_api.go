// Copyright 2013, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tabletmanager

import (
	"time"

	log "github.com/golang/glog"
	"github.com/youtube/vitess/go/vt/mysqlctl"
	"github.com/youtube/vitess/go/vt/topo"
)

// TabletManagerConn defines the interface used to talk to a remote tablet
type TabletManagerConn interface {
	//
	// Various read-only methods
	//

	// Ping will try to ping the remote tablet
	Ping(tablet *topo.TabletInfo, waitTime time.Duration) error

	// GetSchema asks the remote tablet for its database schema
	GetSchema(tablet *topo.TabletInfo, tables []string, includeViews bool, waitTime time.Duration) (*mysqlctl.SchemaDefinition, error)

	// GetPermissions asks the remote tablet for its permissions list
	GetPermissions(tablet *topo.TabletInfo, waitTime time.Duration) (*mysqlctl.Permissions, error)

	//
	// Various read-write methods
	//

	// ChangeType asks the remote tablet to change its type
	ChangeType(tablet *topo.TabletInfo, dbType topo.TabletType, waitTime time.Duration) error

	//
	// Replication related methods
	//

	// SlavePosition returns the tablet's mysql slave position
	SlavePosition(tablet *topo.TabletInfo, waitTime time.Duration) (*mysqlctl.ReplicationPosition, error)

	// WaitSlavePosition asks the tablet to wait until it reaches that
	// position in mysql replication
	WaitSlavePosition(tablet *topo.TabletInfo, replicationPosition *mysqlctl.ReplicationPosition, waitTime time.Duration) (*mysqlctl.ReplicationPosition, error)

	// MasterPosition returns the tablet's master position
	MasterPosition(tablet *topo.TabletInfo, waitTime time.Duration) (*mysqlctl.ReplicationPosition, error)

	// StopSlave stops the mysql replication
	StopSlave(tablet *topo.TabletInfo, waitTime time.Duration) error

	// GetSlaves returns the addresses of the slaves
	GetSlaves(tablet *topo.TabletInfo, waitTime time.Duration) (*SlaveList, error)

	// WaitBlpPosition asks the tablet to wait until it reaches that
	// position in replication
	WaitBlpPosition(tablet *topo.TabletInfo, blpPosition mysqlctl.BlpPosition, waitTime time.Duration) error

	//
	// Reparenting related functions
	//

	// SlaveWasPromoted tells the remote tablet it is now the master
	SlaveWasPromoted(tablet *topo.TabletInfo, waitTime time.Duration) error

	// SlaveWasRestarted tells the remote tablet its master has changed
	SlaveWasRestarted(tablet *topo.TabletInfo, args *SlaveWasRestartedData, waitTime time.Duration) error
}

type TabletManagerConnFactory func(topo.Server) TabletManagerConn

var tabletManagerConnFactories = make(map[string]TabletManagerConnFactory)

func RegisterTabletManagerConnFactory(name string, factory TabletManagerConnFactory) {
	if _, ok := tabletManagerConnFactories[name]; ok {
		log.Fatalf("RegisterTabletManagerConn %s already exists", name)
	}
	tabletManagerConnFactories[name] = factory
}
