import React, { useMemo } from 'react';
import 'chart.js/auto';
import { Doughnut } from 'react-chartjs-2';
import { Box, Card, CardContent, CardHeader, Divider, Typography, useTheme } from '@mui/material';
import BitcoinIcon from '@mui/icons-material/CurrencyBitcoin';
import TokenIcon from '@mui/icons-material/Token';
import OtherIcon from '@mui/icons-material/DeviceUnknown';

export const UtxosByType = (
  {
    data,
    ...props
  }
) => {
  const theme = useTheme();

  const combinedData = useMemo(() => {
    const combinedData = {
      p2pkh: 0,
      token: 0,
      other: 0,
    }
    if (data) {
      // clean up keys
      for (let [key, value] of Object.entries(data)) {
        if (key === "pubkeyhash") {
          key = "p2pkh";
        } else if (key.match(/token/)) {
          key = "token";
        } else {
          key = "other";
        }
        combinedData[key] += Number(value);
      }
    }
    return combinedData;
  }, [data]);

  const chartData = useMemo(() => {
    const dataPoints = [];
    const dataLabels = [];

    if (data) {


      for (const [key, value] of Object.entries(combinedData)) {
        dataPoints.push(value);
        dataLabels.push(key);
      }
    }

    return {
      datasets: [
        {
          data: dataPoints,
          backgroundColor: ['#3F51B5', '#e53935', '#FB8C00'],
          borderWidth: 2,
          borderColor: '#FFFFFF',
          hoverBorderColor: '#FFFFFF'
        }
      ],
      labels: dataLabels
    };
  }, [combinedData]);

  const options = {
    animation: false,
    cutoutPercentage: 80,
    layout: { padding: 0 },
    legend: {
      display: false
    },
    maintainAspectRatio: false,
    responsive: true,
    tooltips: {
      backgroundColor: theme.palette.background.paper,
      bodyFontColor: theme.palette.text.secondary,
      borderColor: theme.palette.divider,
      borderWidth: 1,
      enabled: true,
      footerFontColor: theme.palette.text.secondary,
      intersect: false,
      mode: 'index',
      titleFontColor: theme.palette.text.primary
    }
  };

  const total = combinedData.p2pkh + combinedData.token + combinedData.other;

  const types = [
    {
      title: 'P2PKH',
      value: (100 * combinedData.p2pkh / total).toFixed(1),
      icon: BitcoinIcon,
      color: '#3F51B5'
    },
    {
      title: 'Tokens',
      value: (100 * combinedData.token / total).toFixed(1),
      icon: TokenIcon,
      color: '#E53935'
    },
    {
      title: 'Other',
      value: (100 * combinedData.other / total).toFixed(1),
      icon: OtherIcon,
      color: '#FB8C00'
    }
  ];

  return (
    <Card {...props}>
      <CardHeader
        title="Utxos by type"
        subheader="Coming soon..."
      />
      <Divider />
      <CardContent>
        <Box
          sx={{
            height: 300,
            position: 'relative'
          }}
        >
          <Doughnut
            data={chartData}
            options={options}
            type=""
          />
        </Box>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            pt: 2
          }}
        >
          {types.map(({
            color,
            icon: Icon,
            title,
            value
          }) => (
            <Box
              key={title}
              sx={{
                p: 1,
                textAlign: 'center'
              }}
            >
              <Icon color="action" />
              <Typography
                color="textPrimary"
                variant="body1"
              >
                {title}
              </Typography>
              <Typography
                style={{ color }}
                variant="h4"
              >
                {value}
                %
              </Typography>
            </Box>
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};
