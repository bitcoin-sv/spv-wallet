import React from 'react';

import { Alert, Card, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { DestinationsList } from "../components/destinations";
import PerfectScrollbar from "react-perfect-scrollbar";
import { useQueryList } from "../hooks/use-query-list";

export const Destinations = () => {
  const { items, loading, error, Pagination } = useQueryList({ modelFunction: 'GetDestinations' });

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        Destinations
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
              <DestinationsList items={items}/>
            </PerfectScrollbar>
            <Pagination/>
          </Card>
        </>
      }
    </DashboardLayout>
  );
};
