import React, { useState } from 'react';
import bsv from 'bsv';
import { BuxClient } from "@buxorg/js-buxclient";

export const useUser = () => {

  const [xPrivString, setXPrivString] = useState('');
  const [xPubString, setXPubString] = useState('');
  const [accessKeyString, setAccessKeyString] = useState('');
  const [server, setServer] = useState('');
  const [transportType, setTransportType] = useState('');
  const [adminKeyString, setAdminKey] = useState('');

  let xPriv, xPub, xPubId, accessKey, adminKey, adminId;

  if (xPrivString) {
    xPriv = bsv.HDPrivateKey.fromString(xPrivString);
    xPub = xPriv.hdPublicKey;
    setXPubString(xPub.toString());
    xPubId = bsv.crypto.Hash.sha256(Buffer.from(xPubString)).toString('hex')
  } else if (xPubString) {
    xPriv = null;
    xPub = bsv.HDPublicKey.fromString(xPubString)
    setXPubString(xPub.toString());
    xPubId = bsv.crypto.Hash.sha256(Buffer.from(xPubString)).toString('hex')
  } else if (accessKeyString) {
    xPriv = null;
    setXPubString(null);
    xPub = null;
    accessKey = bsv.PrivateKey.fromString(accessKeyString);
    xPubId = bsv.crypto.Hash.sha256(Buffer.from(accessKey.publicKey.toString())).toString('hex')
  }

  if (adminKeyString) {
    adminKey = bsv.HDPrivateKey.fromString(adminKeyString);
    adminId = bsv.crypto.Hash.sha256(Buffer.from(adminKeyString)).toString('hex')
  }

  let buxClient, buxAdminClient;
  if (server && transportType) {
    buxClient = new BuxClient(server, {
      transportType: transportType,
      xPriv,
      xPub,
      accessKey,
      signRequest: true,
    });
    if (adminKey) {
      buxAdminClient = new BuxClient(server, {
        transportType: transportType,
        xPriv,
        xPub,
        accessKey,
        signRequest: true,
      });
      buxAdminClient.SetAdminKey(adminKey);
    }
  }

  return {
    xPrivString,
    xPriv,
    xPubString,
    xPub,
    xPubId,
    accessKey,
    accessKeyString,
    transportType,
    server,
    adminKey,
    adminId,
    buxClient,
    buxAdminClient,
    setXPrivString,
    setXPubString,
    setAccessKeyString,
    setServer,
    setTransportType,
    setAdminKey,
  };
}
