import React from 'react';

import { Alert, Card, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { TransactionsList } from "../components/transactions";
import PerfectScrollbar from "react-perfect-scrollbar";
import { useQueryList } from "../hooks/use-query-list";

export const Transactions = () => {
  const { items, loading, error, Pagination } = useQueryList({ modelFunction: 'GetTransactions' });

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        Transactions
      </Typography>
      {loading
        ?
        <>Loading...</>
        :
        <>
          {!!error &&
          <Alert severity="error">{error}</Alert>
          }
          <Card>
            <PerfectScrollbar>
              <TransactionsList items={items}/>
            </PerfectScrollbar>
            <Pagination/>
          </Card>
        </>
      }
    </DashboardLayout>
  );
};
