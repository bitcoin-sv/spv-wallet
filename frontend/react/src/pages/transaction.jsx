import bsv from 'bsv';
import React, { useEffect, useState } from 'react';
import { BuxClient } from '@buxorg/js-buxclient';

import { Alert, TextField, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { useUser } from "../hooks/user";
import { useLocation } from "react-router-dom";
import { JsonView } from "../components/json-view";

export const Transaction = () => {
  const { xPriv, xPub, accessKey, server, transportType } = useUser();
  const location = useLocation();
  const params = new URLSearchParams(location.search)

  const [ txId, setTxId ] = useState('');
  const [ transaction, setTransaction ] = useState(null);
  const [ loading, setLoading ] = useState(false);
  const [ error, setError ] = useState('');

  const buxClient = new BuxClient(server, {
    transportType: transportType,
    xPriv,
    xPub,
    accessKey,
    signRequest: true,
  });
  buxClient.SetSignRequest(true);

  useEffect(() => {
    const tx_id = params.get('tx_id');
    if (tx_id) {
      setTxId(tx_id);
    }
  }, [params]);

  useEffect(() => {
    if (txId) {
      setLoading(true);
      buxClient.GetTransaction(txId.trim()).then(tx => {
        setTransaction(tx);
        setError('');
        setLoading(false);
      }).catch(e => {
        setTransaction(null);
        setError(e.message);
        setLoading(false);
      });
    }
  },[txId]);

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        Transaction
      </Typography>
      <TextField
        fullWidth
        label="Transaction ID"
        margin="normal"
        value={txId}
        onChange={(e) => setTxId(e.target.value)}
        type="text"
        variant="outlined"
      />
      {loading
      ?
        <>Loading...</>
      :
        <>
          {!!error &&
          <Alert severity="error">{error}</Alert>
          }
          {txId && transaction && <>
            <h2>Bux transaction</h2>
            <JsonView jsonData={transaction} />
            <h2>Bitcoin transaction</h2>
            <JsonView jsonData={(new bsv.Transaction(transaction.hex)).toObject()} />
          </>}
        </>
      }
    </DashboardLayout>
  );
};
