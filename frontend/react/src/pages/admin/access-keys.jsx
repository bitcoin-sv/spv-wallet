import React, { useEffect, useState } from 'react';

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { XPubsList } from "../../components/xpubs";
import { useNavigate } from "react-router-dom";
import { AccessKeysList } from "../../components/access-keys";
import { Box } from "@mui/material";

export const AdminAccessKeys = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ showRevoked, setShowRevoked ] = useState(false);
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/access-keys');
    }
  }, [params]);

  useEffect(() => {
    const conditions = {}
    if (filter) {
      conditions["$or"] = [
        {
          id: filter
        },
        {
          xpub_id: filter
        },
      ];
    }
    if (!showRevoked) {
      conditions.revoked_at = {$exists: false}
    }
    setConditions(conditions)
  }, [filter, showRevoked]);

  const additionalFilters = function() {
    return (
      <Box style={{ marginLeft: 20 }}>
        Show Revoked <input type="checkbox" checked={showRevoked} onClick={() => setShowRevoked(!showRevoked)}/>
      </Box>
    )
  }

  return (
    <DashboardLayout>
      <AdminListing
        key="admin_access-keys_listing"
        modelFunction="AdminGetAccessKeys"
        title="Access Keys"
        ListingComponent={AccessKeysList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
        additionalFilters={additionalFilters}
      />
    </DashboardLayout>
  );
};
