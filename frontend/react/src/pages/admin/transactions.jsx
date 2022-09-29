import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

import { DashboardLayout } from "../../components/dashboard-layout";
import { AdminListing } from "../../components/listing/admin";
import { TransactionsList } from "../../components/transactions";

export const AdminTransactions = () => {
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search)

  const [ filter, setFilter ] = useState('');
  const [ conditions, setConditions ] = useState(null);

  useEffect(() => {
    const search = params.get('search');
    if (search) {
      setFilter(search);
      navigate('/admin/transactions');
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
        key="admin_transactions_listing"
        modelFunction="AdminGetTransactions"
        title="Transactions"
        ListingComponent={TransactionsList}
        filter={filter}
        setFilter={setFilter}
        conditions={conditions}
      />
    </DashboardLayout>
  );
};
