import bsv from 'bsv';
import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

import { Alert, Button, TextField, Typography } from "@mui/material";

import { DashboardLayout } from "../../components/dashboard-layout";
import { useUser } from "../../hooks/user";

export const AdminRegisterXPub = () => {
  const navigate = useNavigate();
  const { buxAdminClient } = useUser();

  const [ newXPub, setNewXPub ] = useState("");
  const [ xPrivInput, setXPrivInput ] = useState("");
  const [ loading, setLoading ] = useState(false);
  const [ error, setError ] = useState('');

  useEffect(() => {
    if (!buxAdminClient) {
      navigate('/');
    }
  }, [buxAdminClient]);

  useEffect(() => {
    if (xPrivInput) {
      try {
        const xPrivHD = bsv.HDPrivateKey.fromString(xPrivInput); // will throw on error
        setNewXPub(xPrivHD.hdPublicKey.toString());
        setXPrivInput("");
      } catch(e) {
        setError(e.message);
      }
    }
  }, [xPrivInput]);

  const handleRegisterXPub = useCallback((newXPub) => {
    if (!newXPub) {
      setError("No xPub to add")
    }
    setLoading(true);
    try {
      const xPubHD = bsv.HDPublicKey.fromString(newXPub); // will throw on error
      buxAdminClient.RegisterXpub(newXPub);
      alert("XPub added");
      setNewXPub("");
    } catch(e) {
      setError(e.message);
    }
    setLoading(false);
  }, [buxAdminClient]);

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        Register xPub
      </Typography>
      <Typography
        color="error"
        variant="subtitle1"
        style={{
          marginLeft: 20,
          marginRight: 20,
        }}
      >
        Only register xPub's that are new or not registered with another server. If multiple servers are managing an xPub, the xPub state can get out of synch.
      </Typography>
      <TextField
        fullWidth
        label="xPriv"
        margin="normal"
        value={xPrivInput}
        onChange={(e) => setXPrivInput(e.target.value)}
        type="text"
        variant="outlined"
      />
      <Typography
        color="inherit"
        variant="subtitle2"
        style={{
          marginLeft: 20,
          marginRight: 20,
        }}
      >
        Paste an xPriv to convert to xPub. The xPriv will be converted and discarded.
      </Typography>
      <TextField
        fullWidth
        label="xPub"
        margin="normal"
        value={newXPub}
        onChange={(e) => setNewXPub(e.target.value)}
        type="text"
        variant="outlined"
      />
      <Button
        sx={{ mr: 1 }}
        onClick={() => handleRegisterXPub(newXPub)}
      >
        + Register xPub
      </Button>
      {loading
      ?
        <>Loading...</>
      :
        <>
          {!!error &&
            <Alert severity="error">{error}</Alert>
          }
        </>
      }
    </DashboardLayout>
  );
};
