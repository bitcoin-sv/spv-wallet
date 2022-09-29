import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

import { CacheProvider } from '@emotion/react';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { CssBaseline } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';

import { createEmotionCache } from './utils/create-emotion-cache';
import { theme } from './theme';
import { useUser } from "./hooks/user";

import Dashboard from './pages/dashboard';
import Login from "./pages/login";
import { Transactions } from "./pages/transactions";
import NotFound from "./pages/404";
import { Transaction } from "./pages/transaction";
import { XPub } from "./pages/xpub";
import { AccessKeys } from "./pages/access-keys";
import { Destinations } from "./pages/destinations";
import { Destination } from "./pages/destination";
import { DestinationNew } from "./pages/destination-new";
import { TransactionNew } from "./pages/transaction-new";
import { AdminPaymails } from "./pages/admin/paymails";
import { AdminRegisterXPub } from "./pages/admin/register-xpub";
import { AdminXPubs } from "./pages/admin/xpubs";
import { AdminDestinations } from "./pages/admin/destinations";
import { AdminAccessKeys } from "./pages/admin/access-keys";
import { AdminBlockHeaders } from "./pages/admin/block-headers";
import { AdminTransactions } from "./pages/admin/transactions";
import { AdminUtxos } from "./pages/admin/utxos";
import { AdminTransactionRecord } from "./pages/admin/transaction-record";

const clientSideEmotionCache = createEmotionCache();

export const App = () => {
  const { xPub, accessKeyString, adminKey } = useUser();

  return (
    <CacheProvider value={clientSideEmotionCache}>
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <ThemeProvider theme={theme}>
          <CssBaseline/>
          {(!!xPub || !!accessKeyString || !!adminKey)
            ?
            <BrowserRouter>
              <Routes>
                <Route exact path="/xpub" name="xPub" element={<XPub/>} />
                <Route exact path="/destination" name="Destination search" element={<Destination/>} />
                <Route exact path="/destinations" name="Destinations" element={<Destinations/>} />
                <Route exact path="/destination-new" name="New Destinations" element={<DestinationNew/>} />
                <Route exact path="/transaction" name="Transaction search" element={<Transaction/>} />
                <Route exact path="/transactions" name="Transactions" element={<Transactions/>} />
                <Route exact path="/transaction-new" name="New Transactions" element={<TransactionNew/>} />
                <Route exact path="/access-keys" name="Access Keys" element={<AccessKeys/>} />
                {(!!xPub || !!accessKeyString)
                  ?
                  <Route exact path="/" name="xPub" element={<XPub/>}/>
                  :
                  <>
                    <Route exact path="/" name="Admin dashboard" element={<Dashboard/>}/>
                  </>
                }
                <Route exact path="/admin/dashboard" name="Admin dashboard" element={<Dashboard/>}/>
                <Route exact path="/admin/register-xpub" name="Admin register xpub" element={<AdminRegisterXPub/>}/>
                <Route exact path="/admin/access-keys" name="Admin access keys" element={<AdminAccessKeys/>}/>
                <Route exact path="/admin/block-headers" name="Admin block headers" element={<AdminBlockHeaders/>}/>
                <Route exact path="/admin/destinations" name="Admin destinations" element={<AdminDestinations/>}/>
                <Route exact path="/admin/paymails" name="Admin paymails" element={<AdminPaymails/>}/>
                <Route exact path="/admin/transaction-record" name="Admin record transactions" element={<AdminTransactionRecord/>}/>
                <Route exact path="/admin/transactions" name="Admin transactions" element={<AdminTransactions/>}/>
                <Route exact path="/admin/utxos" name="Admin utxos" element={<AdminUtxos/>}/>
                <Route exact path="/admin/xpubs" name="Admin xPubs" element={<AdminXPubs/>}/>
                <Route path="/" name="404" element={<NotFound/>} />
              </Routes>
            </BrowserRouter>
            :
            <Login/>
          }
        </ThemeProvider>
      </LocalizationProvider>
    </CacheProvider>
  );
}
