import React, { useMemo } from 'react';
import 'chart.js/auto';
import { Bar } from 'react-chartjs-2';
import Moment from 'moment';
import { extendMoment } from 'moment-range';

import { Box, Card, CardContent, CardHeader, Divider, useTheme } from '@mui/material';

const moment = extendMoment(Moment);

export const AdminChart = (
  {
    data,
    ...props
  }
) => {
  const theme = useTheme();

  const chartData = useMemo(() => {
    const dataPoints = [];
    const dataLabels = [];

    if (data) {
      const dataKeys = Object.keys(data).map(k => Number(k));
      if (dataKeys.length > 0) {
        const minDate = moment.max(moment().subtract(6, 'months'), moment(Math.min(...dataKeys), 'YYYYMMDD'));
        const maxDate = moment.max(moment().subtract(3, 'months'), moment(Math.min(...dataKeys), 'YYYYMMDD'));
        //const maxDate = moment(Math.max(...dataKeys), 'YYYYMMDD');
        const range = moment.range(minDate, maxDate);
        for (let day of range.by('day')) {
          const d = day.format('YYYYMMDD');
          dataPoints.push(data[d.toString()]);
          dataLabels.push(day.format('YYYY-MM-DD'))
        }
      }
    }

    return {
      datasets: [
        {
          backgroundColor: '#3F51B5',
          barPercentage: 0.5,
          barThickness: 12,
          borderRadius: 4,
          categoryPercentage: 0.5,
          data: dataPoints,
          label: 'Transactions',
          maxBarThickness: 10
        },
      ],
      labels: dataLabels
    };
  }, [data]);

  const options = {
    animation: false,
    cornerRadius: 20,
    layout: { padding: 0 },
    legend: { display: false },
    maintainAspectRatio: false,
    responsive: true,
    xAxes: [
      {
        ticks: {
          fontColor: theme.palette.text.secondary
        },
        gridLines: {
          display: false,
          drawBorder: false
        }
      }
    ],
    yAxes: [
      {
        ticks: {
          fontColor: theme.palette.text.secondary,
          beginAtZero: true,
          min: 0
        },
        gridLines: {
          borderDash: [2],
          borderDashOffset: [2],
          color: theme.palette.divider,
          drawBorder: false,
          zeroLineBorderDash: [2],
          zeroLineBorderDashOffset: [2],
          zeroLineColor: theme.palette.divider
        }
      }
    ],
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

  return (
    <Card {...props}>
      <CardHeader
        title="Transactions per day"
        subheader="Last 3 months"
      />
      <Divider />
      <CardContent>
        <Box
          sx={{
            height: 400,
            position: 'relative'
          }}
        >
          <Bar
            data={chartData}
            options={options}
          />
        </Box>
      </CardContent>
    </Card>
  );
};
