import React, { useEffect, useState } from 'react';

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { XPubsList } from "../../components/xpubs";
import { useNavigate } from "react-router-dom";
import { DestinationsList } from "../../components/destinations";

export const AdminDestinations = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/destinations');
    }
  }, [params]);

  useEffect(() => {
    if (filter) {
      const conditions = {
        $or: [
          {
            id: filter
          },
          {
            address: filter
          },
          {
            locking_script: filter
          },
          {
            xpub_id: filter
          },
        ],
      }
      setConditions(conditions)
    } else {
      setConditions(null);
    }
  }, [filter]);

  return (
    <DashboardLayout>
      <AdminListing
        key="admin_destinations_listing"
        modelFunction="AdminGetDestinations"
        title="Destinations"
        ListingComponent={DestinationsList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
      />
    </DashboardLayout>
  );
};
