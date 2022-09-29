import { useUser } from "./user";
import React, { useEffect, useState } from "react";
import { TablePagination } from "@mui/material";

export const useQueryList = function (
  {
    modelFunction,
    admin = false,
    conditions,
  }
) {
  const {buxClient, buxAdminClient} = useUser();

  const [items, setItems] = useState([]);
  const [itemsCount, setItemsCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [refreshData, setRefreshData ] = useState(0);

  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);

  const handleLimitChange = (event) => {
    setLimit(event.target.value);
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage);
  };

  useEffect(() => {
    setPage(0);
  }, [limit]);

  useEffect(() => {
    const client = admin ? buxAdminClient : buxClient;
    if (!client) return;
    client[`${modelFunction}Count`](conditions || {}, {}).then(count => {
      setItemsCount(count);
      setError('');
    }).catch(e => {
      setItemsCount(limit);
      setError(e.message);
    });
  }, [refreshData, conditions]);

  useEffect(() => {
    setLoading(true)
    const queryParams = {
      page: page + 1,
      page_size: limit,
      order_by_field: 'created_at',
      sort_direction: 'desc',
    };
    const client = admin ? buxAdminClient : buxClient;
    if (!client) return;
    client[`${modelFunction}`](conditions || {}, {}, queryParams).then(items => {
      setItems([...items]);
      setError('');
      setLoading(false);
    }).catch(e => {
      setError(e.message);
      setLoading(false);
    });
  }, [refreshData, conditions, page, limit]);

  const Pagination = () => {
    return (
      <TablePagination
        component="div"
        count={itemsCount}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleLimitChange}
        page={page}
        rowsPerPage={limit}
        rowsPerPageOptions={[5, 10, 25, 50, 100]}
        showFirstButton={true}
        showLastButton={true}
      />
    )
  };

  return {
    items,
    loading,
    error,
    setError,
    Pagination,
    refreshData,
    setRefreshData,
    buxClient,
    buxAdminClient,
  }
}
