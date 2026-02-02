package controllers

import (
	"fmt"
	"strings"
)

// GenerateResultsEmailHTML creates a rich HTML email with results and audit logs
func GenerateResultsEmailHTML(electionName, winnerName string, candidates []map[string]interface{}, logs []AuditLog) string {

	// 1. Build Candidates Table Rows
	var candRows string
	winnerVotes := 0

	for _, c := range candidates {
		name, _ := c["name"].(string)
		votes := 0
		if v, ok := c["voteCount"]; ok {
			// handle various number types
			switch val := v.(type) {
			case int:
				votes = val
			case int32:
				votes = int(val)
			case int64:
				votes = int(val)
			case float64:
				votes = int(val)
			}
		} else if v, ok := c["votes"]; ok {
			// fallback check
			switch val := v.(type) {
			case int:
				votes = val
			case float64:
				votes = int(val)
			}
		}

		if name == winnerName {
			winnerVotes = votes
		}

		row := fmt.Sprintf(`
			<tr>
				<td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
				<td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold;">%d</td>
			</tr>
		`, name, votes)
		candRows += row
	}

	// 2. Build Audit Log Table Rows
	var logRows string
	if len(logs) == 0 {
		logRows = `<tr><td colspan="3" style="padding:10px; color:#888;">No logs recorded.</td></tr>`
	} else {
		for _, l := range logs {
			timeStr := l.Timestamp.Format("2006-01-02 15:04:05")
			row := fmt.Sprintf(`
				<tr>
					<td style="padding: 8px; border-bottom: 1px solid #f0f0f0; font-size: 12px; color: #666;">%s</td>
					<td style="padding: 8px; border-bottom: 1px solid #f0f0f0; font-size: 13px;">%s</td>
					<td style="padding: 8px; border-bottom: 1px solid #f0f0f0; font-size: 12px; color: #555;">%s</td>
				</tr>
			`, timeStr, l.Action, l.Details)
			logRows += row
		}
	}

	// 3. Assemble Full HTML
	// Utilizing a clean, responsive card design
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f9; color: #333; margin: 0; padding: 0; }
			.container { max-width: 600px; margin: 20px auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
			.header { background: #6c5ce7; color: #ffffff; padding: 30px 20px; text-align: center; }
			.header h2 { margin: 0; font-size: 24px; }
			.content { padding: 30px 20px; }
			.winner-box { background: #e0f7fa; border-left: 5px solid #00bcd4; padding: 15px; margin-bottom: 25px; border-radius: 4px; }
			.winner-title { font-size: 14px; text-transform: uppercase; color: #008ba3; letter-spacing: 1px; font-weight: bold; }
			.winner-name { font-size: 20px; color: #333; margin-top: 5px; }
			.section-title { font-size: 18px; font-weight: 600; margin-bottom: 15px; border-bottom: 2px solid #eee; padding-bottom: 5px; color: #2d3436; }
			table { width: 100%%; border-collapse: collapse; margin-bottom: 25px; }
			th { text-align: left; background: #f8f9fa; padding: 10px; font-size: 12px; text-transform: uppercase; color: #636e72; }
			.footer { background: #f1f2f6; text-align: center; padding: 20px; font-size: 12px; color: #b2bec3; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h2>Election Results</h2>
				<p>%s</p>
			</div>
			
			<div class="content">
				<div class="winner-box">
					<div class="winner-title">Winner Declared</div>
					<div class="winner-name">%s (%d Votes)</div>
				</div>

				<div class="section-title">Vote Count Summary</div>
				<table>
					<thead><tr><th>Candidate</th><th>Votes</th></tr></thead>
					<tbody>
						%s
					</tbody>
				</table>

				<div class="section-title">Election Audit Log</div>
				<p style="font-size: 12px; color: #999; margin-bottom: 10px;">Official timeline of election events.</p>
				<table>
					<thead><tr><th>Time (UTC)</th><th>Action</th><th>Details</th></tr></thead>
					<tbody>
						%s
					</tbody>
				</table>
			</div>

			<div class="footer">
				<p>Thank you for participating in the %s.<br>
				This is an automated message from BlockVotes.</p>
			</div>
		</div>
	</body>
	</html>
	</html>
	`, electionName, winnerName, winnerVotes, candRows, logRows, electionName)

	return strings.TrimSpace(html)
}
