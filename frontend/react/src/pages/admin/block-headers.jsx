import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { BlockHeadersList } from "../../components/block-headers";

export const AdminBlockHeaders = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/block-headers');
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
            height: Number(filter),
          }
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
        key="admin_block_headers_listing"
        modelFunction="AdminGetBlockHeaders"
        title="Block Headers"
        ListingComponent={BlockHeadersList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
      />
    </DashboardLayout>
  );
};
