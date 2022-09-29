import React, { useCallback, useState } from 'react';
import { BuxClient } from '@buxorg/js-buxclient';
import QRCode from "react-qr-code";

import { Alert, Box, Button, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { useUser } from "../hooks/user";
import { JsonView } from "../components/json-view";

export const DestinationNew = () => {
  const { xPriv, server, transportType } = useUser();

  const [ destination, setDestination ] = useState(null);

  const [ loading, setLoading ] = useState(false);
  const [ error, setError ] = useState('');

  const buxClient = new BuxClient(server, {
    transportType: transportType,
    xPriv,
    signRequest: true,
  });

  const handleNewDestination = useCallback(() => {
    setLoading(true);
    buxClient.NewDestination({}).then(d => {
      setDestination(d);
      setError('');
      setLoading(false);
    }).catch(e => {
      setError(e.message);
      setLoading(false);
    });
  },[]);

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        New Destination
      </Typography>
      <Button
        onClick={() => {
          handleNewDestination();
        }}
      >
        Create new destination
      </Button>
      {loading
      ?
        <>Loading...</>
      :
        <>
          {!!error &&
          <Alert severity="error">{error}</Alert>
          }
          {destination && <>
            <h2>Bux destination</h2>
            <JsonView jsonData={destination} />
            <Box display="flex" justifyContent="center">
              <QRCode value={`bitcoinsv:${destination.address}`} />
            </Box>
          </>}
        </>
      }
    </DashboardLayout>
  );
};
