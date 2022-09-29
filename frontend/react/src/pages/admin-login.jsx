import bsv from "bsv";
import React, { useState } from 'react';

import { Box, Button, Container, TextField, Typography } from '@mui/material';

import { useUser } from "../hooks/user";
import { SeverityPill } from "../components/severity-pill";
import { BuxClient } from "@buxorg/js-buxclient";

const AdminLogin = () => {
  const { server, transportType, xPrivString, xPubString, accessKey, setAdminKey } = useUser();

  const [ xPriv, setXPriv ] = useState('');
  const [ error, setError ] = useState('');

  const handleSubmit = async function(e) {
    e.preventDefault();
    if (xPriv) {
      try {
        // try to make a connection and get the xpub
        const buxClient = new BuxClient(server, {
          transportType: transportType,
          xPrivString: xPrivString,
          xPubString: xPubString,
          accessKeyString: accessKey,
          signRequest: true,
        });
        buxClient.SetAdminKey(xPriv)
        const status = await buxClient.AdminGetStatus();

        const key = bsv.HDPrivateKey.fromString(xPriv);
        setAdminKey(xPriv);
      } catch (e) {
        setError(e.reason || e.message);
      }
    }
  }

  return (
    <Box
      component="main"
      sx={{
        alignItems: 'center',
        display: 'flex',
        flexGrow: 1,
        minHeight: '100%'
      }}
    >
      <Container maxWidth="sm">
        <form onSubmit={handleSubmit}>
          <Box sx={{ my: 3 }}>
            <Typography
              color="textPrimary"
              variant="h4"
            >
              Set admin key
            </Typography>
          </Box>
          <TextField
            fullWidth
            label="xPriv"
            margin="normal"
            value={xPriv}
            onChange={(e) => setXPriv(e.target.value)}
            type="text"
            variant="outlined"
          />
          {!!error &&
            <SeverityPill color={"error"}>
              {error}
            </SeverityPill>
          }
          <Box sx={{ py: 2 }}>
            <Button
              color="primary"
              fullWidth
              size="large"
              type="submit"
              variant="contained"
            >
              Set key
            </Button>
          </Box>
        </form>
      </Container>
    </Box>
  );
};

export default AdminLogin;
