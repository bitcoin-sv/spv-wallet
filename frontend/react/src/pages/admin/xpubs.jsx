import React, { useEffect, useState } from 'react';

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { XPubsList } from "../../components/xpubs";
import { useNavigate } from "react-router-dom";

export const AdminXPubs = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/xpubs');
    }
  }, [params]);

  useEffect(() => {
    if (filter) {
      const conditions = {
        id: filter
      }
      setConditions(conditions)
    } else {
      setConditions(null);
    }
  }, [filter]);

  return (
    <DashboardLayout>
      <AdminListing
        key="admin_xpubs_listing"
        modelFunction="AdminGetXPubs"
        title="XPubs"
        ListingComponent={XPubsList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
      />
    </DashboardLayout>
  );
};
