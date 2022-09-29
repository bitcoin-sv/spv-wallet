import React, { useEffect, useState } from 'react';

import { DashboardLayout } from "../../components/dashboard-layout";
import { PaymailsList } from "../../components/paymails";
import { AdminListing } from "../../components/listing/admin";
import { Box } from "@mui/material";

export const AdminPaymails = () => {

  const [ filter, setFilter ] = useState('');
  const [ showDeleted, setShowDeleted ] = useState(false);
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const conditions = {}
    if (filter) {
      const paymail = filter.split('@');
      conditions.alias = paymail[0]
      if (paymail[1]) {
        conditions.domain = paymail[1];
      }
    }
    if (!showDeleted) {
      conditions.deleted_at = {$exists: false}
    }
    setConditions(conditions)
  }, [filter, showDeleted]);

  const additionalFilters = function() {
    return (
      <Box style={{ marginLeft: 20 }}>
        Show Deleted <input type="checkbox" checked={showDeleted} onClick={() => setShowDeleted(!showDeleted)}/>
      </Box>
    )
  }

  return (
    <DashboardLayout>
      <AdminListing
        key="admin_paymails_listing"
        modelFunction="AdminGetPaymails"
        title="Paymails"
        ListingComponent={PaymailsList}
        setFilter={setFilter}
        conditions={conditions}
        additionalFilters={additionalFilters}
      />
    </DashboardLayout>
  );
};
