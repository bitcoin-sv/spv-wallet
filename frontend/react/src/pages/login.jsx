import bsv from "bsv";
import React, { useState } from 'react';

import { Box, Button, Container, Select, TextField, Typography } from '@mui/material';
import MenuItem from '@mui/material/MenuItem';

import { useUser } from "../hooks/user";
import { SeverityPill } from "../components/severity-pill";
import { useLocalStorage } from "../hooks/localstorage";
import { BuxClient } from "@buxorg/js-buxclient";

const publicSettings = {
  transportType: '',
  serverUrl: '',
  loginTitle: '',
  loginSubtitle: '',
  hideServerUrl: ''
}

const Login = () => {
  const { setAccessKeyString, setAdminKey, setServer, setTransportType, setXPrivString } = useUser();

  const [ loginKey, setLoginKey ] = useState('');
  const [ transport, setTransport ] = useLocalStorage('login.transport', publicSettings.transportType);
  const [ serverUrl, setServerUrl ] = useLocalStorage('login.serverUrl', publicSettings.serverUrl);
  const [ error, setError ] = useState('');

  const handleSubmit = async function(e) {
    e.preventDefault();
    if (loginKey && serverUrl && transport) {
      let buxClient;
      let useTransport = transport;
      let useServerUrl = serverUrl
      try {
        if (publicSettings.transportType && publicSettings.serverUrl) {
          // use the hardcoded defaults for transport and server url
          useTransport = publicSettings.transportType;
          useServerUrl = publicSettings.serverUrl;
        }

        // try to make a connection and get the xpub
        buxClient = new BuxClient(useServerUrl, {
          transportType: useTransport,
          xPrivString: loginKey.match(/^xprv/) ? loginKey : '',
          accessKeyString: loginKey.match(/^[^xp]/) ? loginKey : '',
          signRequest: true,
        });
        const xPub = await buxClient.GetXPub();

        if (loginKey.match(/^xprv/)) {
          const key = bsv.HDPrivateKey.fromString(loginKey);
          setXPrivString(loginKey);
        } else {
          const key = bsv.PrivateKey.fromString(loginKey);
          setAccessKeyString(loginKey);
        }
        setServer(useServerUrl);
        setTransportType(useTransport);
      } catch (e) {
        // check whether this is an admin only login
        try {
          if (loginKey.match(/^xprv/)) {
            buxClient.SetAdminKey(loginKey);
            const admin = await buxClient.AdminGetStatus();
            if (admin === true) {
              setAdminKey(loginKey);
              setServer(useServerUrl);
              setTransportType(useTransport);
            }
          }
          return
        } catch(e) {
          //console.error(e);
        }

        console.error(e);
        setError(e.reason || e.message);
      }
    } else {
      setError("Please set a server and an xPriv to connect");
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
              {publicSettings.loginTitle || "Sign in to a Bux server"}
            </Typography>
            <Typography
              color="textSecondary"
              gutterBottom
              variant="body2"
            >
              {publicSettings.loginSubtitle || "Sign in using your xPriv or access key"}
            </Typography>
          </Box>
          <Box>
            <img
              src="/svg/my-password.svg"
              alt="Login"
              style={{
                width: "100%",
                marginTop: '-35%',
                marginBottom: '-25%',
                zIndex: -1
              }}
            />
          </Box>
          {!!(publicSettings.transportType && publicSettings.serverUrl)
          ?
            (!publicSettings.hideServerUrl &&
              <>
                <TextField
                  fullWidth
                  label={publicSettings.transportType}
                  margin="dense"
                  value={publicSettings.serverUrl}
                  type="text"
                  variant="outlined"
                  disabled={true}
                />
              </>
            )
          :
            <>
              <Select
                fullWidth
                label="Server transport"
                margin="dense"
                value={transport}
                onChange={(e) => setTransport(e.target.value)}
                type="text"
                variant="outlined"
              >
                <MenuItem value="graphql">GraphQL</MenuItem>
                <MenuItem value="http">HTTP</MenuItem>
              </Select>
              <TextField
                fullWidth
                label="Server"
                margin="dense"
                value={serverUrl}
                onChange={(e) => setServerUrl(e.target.value)}
                type="text"
                variant="outlined"
              />
            </>
          }
          <TextField
            fullWidth
            label="xPriv / access key"
            margin="dense"
            value={loginKey}
            onChange={(e) => setLoginKey(e.target.value)}
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
              Sign In Now
            </Button>
          </Box>
          <Typography
            variant="subtitle2"
          >
            No information is stored or sent to our servers. All actions are done in the client and the xPriv is
            only stored temporarily in memory.
          </Typography>
        </form>
      </Container>
    </Box>
  );
};

export default Login;
