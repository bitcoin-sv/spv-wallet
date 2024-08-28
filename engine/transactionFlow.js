frontend.onConfirmTransaction(password, recipient, satoshis){
    frontend.POST("/webbackend/transaction")
    webbackend.createTransaction() {
        userXPriv = GetUserXPriv(password)
        webbackend.tService.CreateTransaction(userPaymail, userXPriv, recipient, satoshis) {
            recipients = [{
                satoshis,
                recipient
            }]
            metadata = {
                "receiver": recipient,
                "sender": userPaymail,
            }
            webbackend.CreateAndFinalizeTransaction(recipients, metadata) {
                draftTx = webbackend.client.DraftToRecipients(recipients, metadata) {
                    outputs=[{
                        "to": recipient[i].To,
                        "satoshis": recipient[i].Satoshis,
                        "op_return": recipient[i].OpReturn
                    }]
                    webbackend.client.createDraftTransaction() {
                        webbackend.client.POST("spv-wallet/transaction", {
                            "config": {
                                "outputs": outputs
                            },
                            "metadata": metadata
                        }) 
                        draftTx = spvwallet.newTransaction(xPub, models.TransactionConfig, metadata) {
                            config: spvwallet.engine.TransactionConfig = MapTransactionConfigEngineToModel(models.TransactionConfig)
                            spvwallet.engine.NewTransaction() {
                                draft = spvwallet.engine.newDraftTransaction(rawXPubKey, config) {
                                    draft = DraftTransaction{
                                        //defaults, config, xPubID, newID, expiresAt(defaultDraftTxExpiresIn)
                                    }
                                    draft.createTransactionHex() {
                                        satoshisNeeded = getTotalSatoshis()
                                        draft.processConfigOutputs() {
                                            paymailFrom = GetPaymailAddressesByXPubID(xpubID, metadata?.sender) //but if doesn't exist, logic proceeds
                                            //if m.Configuration.SendAllTo != nil { but in this case it is nil
                                            for(output in config.outputs) {
                                                output.processOutput(paymailFrom) {
                                                    if(output.To("is handcash or relayx 'Handle'")) {
                                                        output.To = ConvertHandle() // // ConvertHandle will convert a $handle or 1handle to a paymail address
                                                    }

                                                    if(output.To.contains("@")) {
                                                        //treating as paymail transaction 
                                                        output.processPaymailOutput(paymailFrom) {
                                                            alias, domain = paymail.SanitizePaymail(output.To)
                                                            output.PaymailP4 = {
                                                                "alias": alias,
                                                                "domain": domain
                                                            }
                                                            if(hasP2P(getCapabilities(domain))) {
                                                                format = "basic" | "beef"
                                                                output.processPaymailViaP2P(p2pDestinationURL, p2pSubmitTxURL, fromPaymail, format) {
                                                                    //weird hack
                                                                    if output.satoshis <= 0 {
                                                                        output.satoshis = 100
                                                                    }
                                                                    //
                                                                    destinationInfo = startP2PTransaction(t.PaymailP4.Alias, t.PaymailP4.Domain, p2pDestinationURL, satoshis) {
                                                                        paymail.GetP2PPaymentDestination(t.PaymailP4.Alias, t.PaymailP4.Domain, p2pDestinationURL, satoshis) {
                                                                            resp = paymail.POST(replaceAliasDomain(p2pDestinationURL), {"satoshis": satoshis}) {
                                                                                returns({
                                                                                    "outputs": [{
                                                                                        "address": address,
                                                                                        "satoshis": satoshis,
                                                                                        "script": "sth"
                                                                                    }],
                                                                                    "reference": "ref"
                                                                                    //go-paymail.PaymentDestinationPayload {Outputs: []*github.com/bitcoin-sv/go-paymail.PaymentOutput len: 1, cap: 1, [*(*"github.com/bitcoin-sv/go-paymail.PaymentOutput")(0xc001a85740)], Reference: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJ2ZXJzaW9uIjoxLCJkZXJpdmF0aW9uUm9vdFBhdGgiOiJtLzUiLCJkZXJpdmF0aW9uUm9vdEluZGV4Ijo1MSwiY291bnQiOjF9.t0UHzz-VronKzWGDN2JUfT5xKoVa3rshLzPi6P2G_29j3vMNxI2vWk290w9U6c8sbMl7l3GhD4dZlCsijznZkg"}
                                                                                })
                                                                            }
                                                                        }
                                                                    }
                                                                    outputValues = utils.SplitOutputValues(satoshis, len(destinationInfo.Outputs)) // SplitOutputValues splits the satoshis value randomly into nrOfOutputs pieces

                                                                    /*
                                                                    // Loop all received P2P outputs and build scripts
                                                                    for index, out := range destinationInfo.Outputs {
                                                                            t.Scripts = append(
                                                                                t.Scripts,
                                                                                &ScriptOutput{
                                                                                    Address:    out.Address,
                                                                                    Satoshis:   outputValues[index],
                                                                                    Script:     out.Script,
                                                                                    ScriptType: utils.ScriptTypePubKeyHash,
                                                                                },
                                                                            )
                                                                        }

                                                                        // Set the remaining P2P information
                                                                        t.PaymailP4.ReceiveEndpoint = p2pSubmitTxURL
                                                                        t.PaymailP4.ReferenceID = destinationInfo.Reference
                                                                        t.PaymailP4.ResolutionType = ResolutionTypeP2P
                                                                        t.PaymailP4.FromPaymail = fromPaymail
                                                                        t.PaymailP4.Format = format
                                                                    */
                                                                    
                                                                }
                                                            }
                                                        }
                                                    } else {
                                                        //OP_RETURN or Script but it's not relevant for now
                                                    }
                                                }
                                            }
                                            //END of processOutput
                                        }
                                        inputUtxos: []bt.UTXO = draft.prepareUtxos(satoshisNeeded) {
                                            //if m.Configuration.SendAllTo != nil { but it's nil for my example
                                            draft.prepareSeparateUtxos(satoshisNeeded) {
                                                //if m.Configuration.IncludeUtxos != nil {
                                                reserveSatoshis = satoshisNeeded + draft.estimateFee(m.Configuration.FeeUnit, 0)
                                                reservedUtxos = reserveUtxos(draft.ID, reserveSatoshis, feePerByte) {
                                                    //it picks UTXOs from database and reserve some of them until needed satoshis are satisfied
                                                    //in case of not sufficient funds it "unreserve" these UTXOs 
                                                }
                                                inputUtxos = draft.getInputsFromUtxos(reservedUtxos) {
                                                    //it returns an array of bt.UTXO objects
                                                    /*
                                                    &bt.UTXO{
                                                        TxID:           txIDBytes,
                                                        Vout:           utxo.OutputIndex,
                                                        Satoshis:       utxo.Satoshis,
                                                        LockingScript:  lockingScript,
                                                        SequenceNumber: bt.DefaultSequenceNumber,
                                                    }
                                                    */
                                                }

                                            }
                                        }
                                        tx = bt.NewTx().FromUTXOs(inputUtxos)
                                        draft.calculateAndSetFee(satoshisReserved, satoshisNeeded) {
                                           //...()
                                           satoshisChange = satoshisReserved - satoshisNeeded - fee
                                           newFee = draft.setChangeDestination(satoshisChange) {
                                                //TODO
                                           }
                                        }
                                        addOutputsToTx() {
                                            //add the given outputs to the bt.Tx
                                            //implementations for scriptType: "nulldata"(OP_RETURN?) | "pubkeyhash"(transaction) or "non-standard output script"
                                            if(scriptType="pubkeyhash") {
                                                tx.AddP2PKHOutputFromScript(script, satoshis)
                                            }else {
                                                tx.AddOutput(bt.Output{
                                                    LockingScript: s,
                                                    Satoshis:      "nulldata" ? 0 : sc.Satoshis,
                                                })
                                            }
                                        }
                                        validateOutputsInputs()
                                        draft.Hex = tx.String()
                                    }
                                    returns(draft)
                                }
                                draft.Save()
                            }
                        }
                    }
                }
                hex = webbackend.client.FinalizeTransaction(draftTx) {
                    GetSignedHex(draftTx, userXPriv) {
                        //signs all the inputs using the given xPriv key
                    }
                }
                webbackend.tryRecordTransaction(draftID, hex) {
                    webbackend.tryRecord(hex, metadata) {
                        webbackend.client.RecordTransaction(draftID, hex, metadata) {
                            webbackend.client.POST("/transaction/record", {
                                hex, draftID, metadata
                            })
                            spvwallet.engine.RecordTransaction(reqXPub, Hex, ReferenceID("as draftID"), ...engine.WithMetadatas(requestBody.Metadata)) {
                                tx = bt.NewTxFromString(Hex)
                                strategy = new outgoingTx({
                                    "BtTx":           btTx,
                                    "RelatedDraftID": draftID,
                                    "XPubKey":        xPubKey,
                                })
                                recordTransaction(strategy) {
                                    strategy.Execute() {
                                        transaction = _createOutgoingTxToRecord(strategy) {
                                            tx = txFromBtTx(oTx.BtTx)
                                            //// Create NEW transaction model
                                            _hydrateOutgoingWithDraft()
                                            _hydrateOutgoingWithSync() {
                                                //(...)
                                                sync = newSyncTransaction()
                                                sync.SyncStatus = SyncStatusPending // wait until transaction is broadcasted or P2P provider is notified
                                            }
                                            tx.processUtxos() {
                                                tx._processInputs() {
                                                    for(input in parsedTx.Inputs) {
                                                        utxo = getUtxo(input.PreviousTxID(), input.PreviousTxOutIndex)
                                                        if(client.IsIUCEnabled()) { //it is enabled on standard transaction
                                                            // check whether the utxo is spent
                                                            // return error if not isReserved or it is reserved for other transaction
                                                        }
                                                        tx.XpubOutputValue[utxo.XpubID] -= int64(utxo.Satoshis)
                                                        // Mark utxo as spent
                                                    }
                                                }
                                                tx._processOutputs() {
                                                    for(output in parsedIx.Outputs) {
                                                        // only save outputs with a satoshi value attached to it
                                                        lockingScript = output.LockingScript.String() 
                                                        //but for STAS token the lockingScript is returned by GetLockingScriptFromSTASLockingScript
                                                        destination = getDestinationByLockingScript(lockingScript)
                                                        tx.XpubOutputValue[destination.XpubID] += int64(output.Satoshis)
                                                        utxo = GetUtxoByTransactionID(tx.ID) ?? newUtxo(destination.XpubID, tx.ID, txLockingScript)
                                                        m.utxos.append(utxo)
                                                        tx.m.XpubOutIDs.append(destination.XpubID) //if not alreadt exists
                                                    }
                                                }
                                                m.TotalValue, m.Fee = m.getValues()
                                            }

                                        }
                                        transaction.Save()
                                        if (transaction.syncTransaction.P2PStatus == SyncStatusReady) {
                                            _outgoingNotifyP2p(transaction) {
                                                processP2PTransaction(transaction) {
                                                    syncTx = transaction.syncTransaction
                                                    syncResults = _notifyPaymailProviders(transaction) {
                                                        for(output in ransaction.draftTransaction.Configuration.Outputs) {
                                                            if (output.PaymailP4) {
                                                                finalizeP2PTransaction(transaction, outputP4) {
                                                                    p2pTransaction = buildP2pTx(transaction, outputP4) {
                                                                        returns({
                                                                            MetaData: {
                                                                                Note:   p4.Note,
			                                                                    Sender: p4.FromPaymail,
                                                                            },
                                                                            Reference: p4.ReferenceID,
                                                                            Beef: ToBeef(transaction), //if p4.Format is BeefPaymailPayloadFormat
                                                                            Hex: transaction.Hex //if p4.Format is BasicPaymailPayloadFormat
                                                                        })
                                                                    }
                                                                    resp :{note, txid} = paymail.SendP2PTransaction(p4.ReceiveEndpoint, p4.Alias, p4.Domain, p2pTransaction)
                                                                }
                                                            }
                                                        }
                                                    }
                                                    syncTx.Results.Results.append(...syncResults)
                                                    syncTx.P2PStatus = SyncStatusComplete
                                                    syncTx.SyncStatus = SyncStatusReady // if it was "pending"
                                                    syncTx.Save()
                                                }

                                            }
                                        }
                                        syncTx = GetSyncTransactionByID(transaction.ID) // get newest syncTx from DB - if it's an internal tx it could be broadcasted by us already
                                        if(syncTx.BroadcastStatus == SyncStatusReady) {
                                            transaction.syncTransaction = syncTx
                                            _outgoingBroadcast(syncTx) {
                                                broadcastSyncTransaction(syncTx) {
                                                    txHex, hexFormat = _getTxHexInFormat() {
                                                        returns("chainstate supports it and can be sucessfully converterd") ? EfFormat: RawTxFormat
                                                    }
                                                    chainstate.Broadcast(syncTx.ID, txHex, hexFormat) {
                                                        broadcastClient.submitTransaction()
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
                
            }
        }
    }
}
