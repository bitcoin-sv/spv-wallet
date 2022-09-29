import React, { useCallback, useState } from 'react';
import { BuxClient } from '@buxorg/js-buxclient';

import { Alert, Box, Button, TextField, Typography } from "@mui/material";

import { DashboardLayout } from "../components/dashboard-layout";
import { useUser } from "../hooks/user";
import { JsonView } from "../components/json-view";

export const TransactionNew = () => {
  const { xPriv, server, transportType } = useUser();

  const [ recipients, setRecipients ] = useState([{}]);
  const [ draftTransaction, setDraftTransaction ] = useState(null);
  const [ transaction, setTransaction ] = useState(null);

  const [ loading, setLoading ] = useState(false);
  const [ error, setError ] = useState('');

  const buxClient = new BuxClient(server, {
    transportType: transportType,
    xPriv,
    signRequest: true,
  });

  const handleNewTransaction = useCallback((recipients) => {
    setLoading(true);
    buxClient.DraftToRecipients(recipients, { agent: 'Bux client test' }).then(d => {
      setDraftTransaction(d);
      setError('');
      setLoading(false);
    }).catch(e => {
      setError(e.message);
      setLoading(false);
    });
  },[]);

  const handleSendTransaction = useCallback(draftTransaction => {
    setLoading(true);
    let transactionHex;
    try {
      transactionHex = buxClient.FinalizeTransaction(draftTransaction);
    } catch (e) {
      setDraftTransaction(null);
      setTransaction(null);
      setError(e.message);
      setLoading(false);
    }

    if (transactionHex) {
      buxClient.RecordTransaction(transactionHex, draftTransaction.id, {}).then(t => {
        setTransaction(t);
        setError('');
        setLoading(false);
      }).catch(e => {
        setTransaction(null);
        setError(e.message);
        setLoading(false);
      });
    }
  },[]);

  return (
    <DashboardLayout>
      <Typography
        color="inherit"
        variant="h4"
      >
        New Transaction
      </Typography>
      {recipients.map((recipient, index) => {
        return (
          <Box display="flex" flexDirection="row">
            <TextField
              fullWidth
              label="Address / paymail"
              margin="normal"
              value={recipient.to}
              onChange={(e) => {
                const newRecipients = [...recipients];
                newRecipients[index].to = e.target.value;
                setRecipients(newRecipients)
              }}
              type="text"
              variant="outlined"
            />
            <TextField
              fullWidth
              type="number"
              label="Satoshis"
              margin="normal"
              value={recipient.satoshis}
              onChange={(e) => {
                const newRecipients = [...recipients];
                newRecipients[index].satoshis = Number(e.target.value);
                setRecipients(newRecipients)
              }}
              type="text"
              variant="outlined"
            />
          </Box>
        )
      })}
      <Button
        fullWidth={false}
        onClick={() => {
          setRecipients([...recipients, {}]);
        }}
      >
        Add recipient
      </Button>
      <Button
        onClick={() => {
          handleNewTransaction(recipients);
        }}
      >
        Create draft transaction
      </Button>
      {loading
      ?
        <>Loading...</>
      :
        <>
          {!!error &&
          <Alert severity="error">{error}</Alert>
          }
          {draftTransaction && <>
            <h2>Bux draft transaction</h2>
            <JsonView jsonData={draftTransaction} />
            {transaction
            ?
              <JsonView jsonData={transaction} />
            :
              <Button
                onClick={() => {
                  handleSendTransaction(draftTransaction);
                }}
              >
                Sign and send transaction
              </Button>
            }
          </>}
        </>
      }
    </DashboardLayout>
  );
};
