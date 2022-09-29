import React, { useEffect, useState } from 'react';
import { BuxClient } from '@buxorg/js-buxclient';

import { Alert, Box, TextField, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { useUser } from "../hooks/user";
import { useLocation } from "react-router-dom";
import QRCode from "react-qr-code";
import { JsonView } from "../components/json-view";

export const Destination = () => {
  const { xPriv, xPub, accessKey, server, transportType } = useUser();
  const location = useLocation();
  const params = new URLSearchParams(location.search)

  const [ id, setId ] = useState('');
  const [ address, setAddress ] = useState('');
  const [ lockingScript, setLockingScript ] = useState('');

  const [ destination, setDestination ] = useState(null);

  const [ loading, setLoading ] = useState(false);
  const [ error, setError ] = useState('');

  const buxClient = new BuxClient(server, {
    transportType: transportType,
    xPriv,
    xPub,
    accessKey,
    signRequest: true,
  });

  useEffect(() => {
    const id = params.get('id');
    if (id) {
      setId(id);
    }
  }, [params]);

  useEffect(() => {
    const address = params.get('address');
    if (address) {
      setAddress(address);
    }
  }, [params]);

  useEffect(() => {
    const lockingScript = params.get('lockingScript');
    if (lockingScript) {
      setLockingScript(lockingScript);
    }
  }, [params]);

  useEffect(() => {
    if (id) {
      setLoading(true);
      buxClient.GetDestinationByID(id).then(d => {
        setDestination(d);
        setError('');
        setLoading(false);
      }).catch(e => {
        setDestination(null);
        setError(e.message);
        setLoading(false);
      });
    } else if (address) {
      setLoading(true);
      buxClient.GetDestinationByAddress(address).then(d => {
        setDestination(d);
        setError('');
        setLoading(false);
      }).catch(e => {
        setDestination(null);
        setError(e.message);
        setLoading(false);
      });
    } else if (lockingScript) {
      setLoading(true);
      buxClient.GetDestinationByLockingScript(lockingScript).then(d => {
        setDestination(d);
        setError('');
        setLoading(false);
      }).catch(e => {
        setDestination(null);
        setError(e.message);
        setLoading(false);
      });
    }
  },[id, address, lockingScript]);

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        Destination
      </Typography>
      <Box display="flex" flexDirection="row" width="100%">
        <TextField
          fullWidth
          label="Destination ID"
          margin="normal"
          value={id}
          onChange={(e) => setId(e.target.value)}
          type="text"
          variant="outlined"
        />
        <TextField
          fullWidth
          label="Destination address"
          margin="normal"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          type="text"
          variant="outlined"
        />
        <TextField
          fullWidth
          label="Destination locking script"
          margin="normal"
          value={lockingScript}
          onChange={(e) => setLockingScript(e.target.value)}
          type="text"
          variant="outlined"
        />
      </Box>
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
            <h2>Transactions</h2>
            <i>Not possible yet ...</i>
          </>}
        </>
      }
    </DashboardLayout>
  );
};
