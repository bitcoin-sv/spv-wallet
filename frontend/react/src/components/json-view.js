import React, { useState } from 'react';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Box,
  Card,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@mui/material';
import ArrowRightIcon from '@mui/icons-material/ArrowRight';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles((theme) => ({
  keyCell: {
    display: 'flex',
    flexBasis: 'flex-start',
    whiteSpace: 'nowrap',
    maxWidth: '120px',
    overflow: 'hidden'
  },
  dataCell: {
    whiteSpace: 'wrap',
    wordBreak: 'break-all',
    verticalAlign: 'top',
    minWidth: '120px',
    maxWidth: '320px',
    overflow: 'hidden',
  }
}))

export const JsonView = (
  {
    jsonData,
  }
) => {
  const classes = useStyles();

  const [ showJson, setShowJson ] = useState();

  return (<>
    <Card>
      <PerfectScrollbar>
        <Box style={{ maxWidth: '100vw', overflow: 'auto' }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell className={classes.keyCell}>
                  Key
                </TableCell>
                <TableCell className={classes.dataCell}>
                  Value
                </TableCell>
              </TableRow>
            </TableHead>
            <JsonViewBody jsonData={jsonData}/>
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
    <div>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          alignItems: 'center',
          cursor: 'pointer'
        }}
        onClick={() => {
          setShowJson(!showJson);
        }}
      >
        JSON
        {showJson
          ?
          <ArrowDropDownIcon fontSize={"small"} />
          :
          <ArrowRightIcon fontSize={"small"} />
        }
      </div>
      {showJson &&
        <pre style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}>
          {JSON.stringify(jsonData, null, 4)}
        </pre>
      }
    </div>
  </>);
}

const JsonViewBody = function(
  {
    jsonData,
  }
) {
  const classes = useStyles();

  return (
    <TableBody>
      {Object.keys(jsonData).map(key => (
        <TableRow
          hover
          key={key}
        >
          <TableCell className={classes.keyCell}>
            <b>{key}</b>
          </TableCell>
          <TableCell className={classes.dataCell}>
            {(jsonData[key] && typeof jsonData[key] === 'object')
              ?
              <Table size="small">
                <JsonViewBody jsonData={jsonData[key]} />
              </Table>
              :
              jsonData[key]
            }
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  )
};
