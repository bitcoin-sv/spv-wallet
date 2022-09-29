import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { format } from 'date-fns';
import { Table, TableBody, TableCell, TableHead, TableRow } from '@mui/material';
import { Link } from "react-router-dom";
import { JsonView } from "./json-view";

export const TransactionsList = ({items}) => {
  const [selectedTransactions, setSelectedTransactions] = useState([]);

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>ID</TableCell>
          <TableCell>Value</TableCell>
          <TableCell>Date</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {items?.map((transaction) => (
          <React.Fragment key={transaction.id}>
            <TableRow
              hover
              selected={selectedTransactions.indexOf(transaction.id) !== -1}
              onClick={() => {
                if (selectedTransactions.indexOf(transaction.id) !== -1) {
                  setSelectedTransactions([])
                } else {
                  setSelectedTransactions([transaction.id])
                }
              }}
            >
              <TableCell>{transaction.id}</TableCell>
              <TableCell>{transaction.output_value}</TableCell>
              <TableCell>
                {format(new Date(transaction.created_at), 'dd/MM/yyyy hh:mm')}
              </TableCell>
            </TableRow>
            {selectedTransactions.indexOf(transaction.id) !== -1 &&
            <TableRow>
              <TableCell colSpan={3}>
                <Link to={`/transaction?tx_id=${transaction.id}`}>Open Transaction details</Link>
                <JsonView jsonData={transaction}/>
              </TableCell>
            </TableRow>
            }
          </React.Fragment>
        ))}
      </TableBody>
    </Table>
  );
};

TransactionsList.propTypes = {
  items: PropTypes.array.isRequired
};
