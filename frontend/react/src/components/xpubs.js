import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@mui/material';

import { format } from "date-fns";
import { JsonView } from "./json-view";

export const XPubsList = (
  {
    items,
    refetch,
  }
) => {

  const [selectedXPubs, setSelectedXPubs] = useState([]);

  useEffect(() => {
    if (items && items.length === 1) {
      setSelectedXPubs([items[0].id]);
    }
  }, [items]);

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>ID</TableCell>
          <TableCell>Balance</TableCell>
          <TableCell>Created</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {items?.map(xpub => (
          <React.Fragment key={`paymail_${xpub.id}`}>
            <TableRow
              hover
              selected={selectedXPubs.indexOf(xpub.id) !== -1}
              style={{
                opacity: xpub.deleted_at ? 0.5 : 1
              }}
              onClick={() => {
                if (selectedXPubs.indexOf(xpub.id) !== -1) {
                  setSelectedXPubs([])
                } else {
                  setSelectedXPubs([xpub.id])
                }
              }}
            >
              <TableCell>{xpub.id}</TableCell>
              <TableCell>{xpub.current_balance}</TableCell>
              <TableCell>
                {format(new Date(xpub.created_at), 'dd/MM/yyyy hh:mm')}
              </TableCell>
            </TableRow>
            {selectedXPubs.indexOf(xpub.id) !== -1 &&
              <TableRow>
                <TableCell colSpan={5}>
                  <JsonView jsonData={xpub}/>
                </TableCell>
              </TableRow>
            }
          </React.Fragment>
        ))}
      </TableBody>
    </Table>
  );
};

XPubsList.propTypes = {
  items: PropTypes.array.isRequired,
  refetch: PropTypes.func,
};
