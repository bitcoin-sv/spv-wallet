package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
)

// Cron job names to be used in WithCronCustomPeriod
const (
	CronJobNameDraftTransactionCleanUp  = "draft_transaction_clean_up"
	CronJobNameSyncTransactionBroadcast = "sync_transaction_broadcast"
	CronJobNameSyncTransactionSync      = "sync_transaction_sync"
	CronJobNameCalculateMetrics         = "calculate_metrics"
)

type cronJobHandler func(ctx context.Context, client *Client) error

// here is where we define all the cron jobs for the client
func (c *Client) cronJobs() taskmanager.CronJobs {
	jobs := taskmanager.CronJobs{}

	addJob := func(name string, period time.Duration, task cronJobHandler) {
		// handler adds the client pointer to the cronJobTask by using a closure
		handler := func(ctx context.Context) (err error) {
			if metrics, enabled := c.Metrics(); enabled {
				end := metrics.TrackCron(name)
				defer func() {
					success := err == nil
					end(success)
				}()
			}
			err = task(ctx, c)
			return
		}

		jobs[name] = taskmanager.CronJob{
			Handler: handler,
			Period:  period,
		}
	}

	addJob(
		CronJobNameDraftTransactionCleanUp,
		60*time.Second,
		taskCleanupDraftTransactions,
	)
	addJob(
		CronJobNameSyncTransactionSync,
		5*time.Minute,
		taskSyncTransactions,
	)

	if _, enabled := c.Metrics(); enabled {
		addJob(
			CronJobNameCalculateMetrics,
			15*time.Second,
			taskCalculateMetrics,
		)
	}

	return jobs
}
