import React, { useEffect, useState } from 'react';
import { Box, Container, Grid } from '@mui/material';

import ViewListIcon from "@mui/icons-material/ViewList";
import BitcoinIcon from '@mui/icons-material/CurrencyBitcoin';
import LocationSearchingIcon from '@mui/icons-material/LocationSearching';
import PaymailIcon from '@mui/icons-material/Message';

import { AdminChart } from '../components/dashboard/chart';
import { DashboardLayout } from '../components/dashboard-layout';
import { useUser } from "../hooks/user";
import AdminLogin from "./admin-login";
import { AdminCard } from "../components/dashboard/card";
import { UtxosByType } from "../components/dashboard/utxos-by-type";

const Dashboard = () => {
  const { buxAdminClient, adminId } = useUser();

  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!buxAdminClient) return

    setLoading(true)
    buxAdminClient.AdminGetStats().then(stats => {
      setStats(stats);
      setError('');
      setLoading(false);
    }).catch(e => {
      setError(e.message);
      setLoading(false);
    });
  }, [adminId]);

  if (!buxAdminClient) {
    return (
      <DashboardLayout>
        <AdminLogin/>
      </DashboardLayout>
    );
  }

  console.log(stats)

  return (
    <DashboardLayout>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8
        }}
      >
        {loading
          ?
          <>Loading...</>
          :
          <Container maxWidth={false}>
            <Grid
              container
              spacing={3}
            >
              <Grid
                item
                lg={3}
                sm={6}
                xl={3}
                xs={12}
              >
                <AdminCard
                  Icon={BitcoinIcon}
                  iconColor={'blue'}
                  title="XPubs"
                  value={stats?.xpubs ? new Intl.NumberFormat().format(stats.xpubs) : ''}
                  listLink="/admin/xpubs"
                />
              </Grid>
              <Grid
                item
                xl={3}
                lg={3}
                sm={6}
                xs={12}
              >
                <AdminCard
                  Icon={ViewListIcon}
                  iconColor={'red'}
                  title="Transactions"
                  value={stats?.transactions ? new Intl.NumberFormat().format(stats.transactions) : ''}
                  listLink="/admin/transactions"
                />
              </Grid>
              <Grid
                item
                xl={3}
                lg={3}
                sm={6}
                xs={12}
              >
                <AdminCard
                  Icon={LocationSearchingIcon}
                  iconColor={'green'}
                  title="Destinations"
                  value={stats?.destinations ? new Intl.NumberFormat().format(stats.destinations) : ''}
                  listLink="/admin/destinations"
                />
              </Grid>
              <Grid
                item
                xl={3}
                lg={3}
                sm={6}
                xs={12}
              >
                <AdminCard
                  Icon={PaymailIcon}
                  iconColor={'grey'}
                  title="Paymails"
                  value={stats?.paymail_addresses ? new Intl.NumberFormat().format(stats.paymail_addresses) : ''}
                  listLink="/admin/paymails"
                />
              </Grid>
              <Grid
                item
                lg={8}
                md={12}
                xl={9}
                xs={12}
              >
                <AdminChart
                  data={stats?.transactions_per_day}
                />
              </Grid>
              <Grid
                item
                lg={4}
                md={6}
                xl={3}
                xs={12}
              >
                <UtxosByType sx={{height: '100%'}} data={stats?.utxos_per_type}/>
              </Grid>
            </Grid>
          </Container>
        }
      </Box>
    </DashboardLayout>
  );
}

export default Dashboard;
