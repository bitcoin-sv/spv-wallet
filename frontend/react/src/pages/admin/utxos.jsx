import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { UtxosList } from "../../components/utxos";
import { Box } from "@mui/material";

export const AdminUtxos = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ showSpent, setShowSpent ] = useState(false);
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/utxos');
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
        {
          script_pub_key: filter
        },
        {
          type: filter
        },
      ];
    }
    if (!showSpent) {
      conditions.spending_tx_id = {$exists: false}
    }
    setConditions(conditions)
  }, [filter, showSpent]);

  const additionalFilters = function() {
    return (
      <Box style={{ marginLeft: 20 }}>
        Show Spent <input type="checkbox" checked={showSpent} onClick={() => setShowSpent(!showSpent)}/>
      </Box>
    )
  }

  return (
    <DashboardLayout>
      <AdminListing
        key="admin_utxos_listing"
        modelFunction="AdminGetUtxos"
        title="Utxos"
        ListingComponent={UtxosList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
        additionalFilters={additionalFilters}
      />
    </DashboardLayout>
  );
};
