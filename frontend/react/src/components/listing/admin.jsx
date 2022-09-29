import React, { useEffect, useState } from "react";
import PropTypes from 'prop-types'
import { useNavigate } from "react-router-dom";
import PerfectScrollbar from "react-perfect-scrollbar";

import { Alert, Box, Card, TextField, Typography } from "@mui/material";

import { useQueryList } from "../../hooks/use-query-list";
import { useDebounce } from "../../hooks/debounce";

export const AdminListing = function(
  {
    modelFunction,
    title,
    ListingComponent,
    conditions,
    filter: initialFilter,
    setFilter,
    additionalFilters,
  }
) {
  const navigate = useNavigate();

  const [ searchFilter, setSearchFilter ] = useState('');
  const debouncedFilter = useDebounce(searchFilter, 500);

  const {
    items,
    loading,
    error,
    Pagination,
    setRefreshData,
    buxAdminClient,
  } = useQueryList({ modelFunction, admin: true, conditions });

  useEffect(() => {
    if (!buxAdminClient) {
      navigate('/');
    }
  }, [buxAdminClient]);

  useEffect(() => {
    if (initialFilter) {
      setSearchFilter(initialFilter);
    }
  }, [initialFilter]);

  useEffect(() => {
    if (setFilter) {
      setFilter(debouncedFilter);
    }
  }, [setFilter, debouncedFilter]);

  return (
    <>
      <Box display="flex" flexDirection="row" alignItems="center">
        <Typography
          color="inherit"
          variant="h4"
        >
          {title}
        </Typography>
        <Box display="flex" flex={1} flexDirection="row" alignItems="center">
          <TextField
            fullWidth
            label="Filter"
            margin="normal"
            value={searchFilter}
            onChange={(e) => setSearchFilter(e.target.value)}
            type="text"
            variant="outlined"
            style={{
              marginLeft: 20
            }}
          />
          {additionalFilters ? additionalFilters() : ''}
        </Box>
      </Box>
      {loading
        ?
        <>Loading...</>
        :
        <>
          {!!error &&
          <Alert severity="error">{error}</Alert>
          }
          <Card>
            <PerfectScrollbar>
              <ListingComponent
                items={items}
                refetch={() => setRefreshData(+new Date())}
              />
            </PerfectScrollbar>
            <Pagination/>
          </Card>
        </>
      }
    </>
  );
}

AdminListing.propTypes = {
  ListingComponent: PropTypes.func.isRequired,
  modelFunction: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  conditions: PropTypes.object,
  filter: PropTypes.string,
  setFilter: PropTypes.func,
  additionalFilters: PropTypes.func,
}
