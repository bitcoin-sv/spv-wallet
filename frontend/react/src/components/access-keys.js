import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { format } from 'date-fns';
import { Button, Table, TableBody, TableCell, TableHead, TableRow, } from '@mui/material';
import { JsonView } from "./json-view";

export const AccessKeysList = ({items, handleRevokeAccessKey}) => {
  const [selectedAccessKeys, setSelectedAccessKeys] = useState([]);

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>ID</TableCell>
          <TableCell>Created</TableCell>
          <TableCell>Revoke</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {items?.map(accessKey => (
          <>
            <TableRow
              hover
              key={`access_key_${accessKey.id}`}
              selected={selectedAccessKeys.indexOf(accessKey.id) !== -1}
              style={{
                opacity: accessKey.revoked_at ? 0.5 : 1
              }}
              onClick={() => {
                if (selectedAccessKeys.indexOf(accessKey.id) !== -1) {
                  setSelectedAccessKeys([])
                } else {
                  setSelectedAccessKeys([accessKey.id])
                }
              }}
            >
              <TableCell>{accessKey.id}</TableCell>
              <TableCell>
                {format(new Date(accessKey.created_at), 'dd/MM/yyyy hh:mm')}
              </TableCell>
              <TableCell>
                {accessKey.revoked_at
                  ?
                  <span title={`Revoked at ${accessKey.revoked_at}`}>Revoked</span>
                  :
                  <Button
                    onClick={() => handleRevokeAccessKey(accessKey)}
                  >
                    Revoke key
                  </Button>
                }
              </TableCell>
            </TableRow>
            {selectedAccessKeys.indexOf(accessKey.id) !== -1 &&
            <TableRow>
              <TableCell colSpan={5}>
                <JsonView jsonData={accessKey}/>
              </TableCell>
            </TableRow>
            }
          </>
        ))}
      </TableBody>
    </Table>
  );
};

AccessKeysList.propTypes = {
  items: PropTypes.array.isRequired,
  handleRevokeAccessKey: PropTypes.func,
};
