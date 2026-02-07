package controllers

import (
	"fmt"
)

// BaseEmailLayout provides a consistent, responsive wrapper for all emails
func BaseEmailLayout(subject, content string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body { margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f9; color: #333; line-height: 1.6; }
        .wrapper { width: 100%%; table-layout: fixed; background-color: #f4f4f9; padding-bottom: 40px; }
        .container { max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #6c5ce7 0%%, #a29bfe 100%%); color: #ffffff; padding: 30px 20px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; }
        .content { padding: 40px 30px; }
        .footer { background: #f1f2f6; text-align: center; padding: 20px; font-size: 12px; color: #b2bec3; border-top: 1px solid #eaeaea; }
        .btn { display: inline-block; background-color: #6c5ce7; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold; margin-top: 20px; }
        .info-box { background: #f8f9fa; padding: 15px; border-left: 4px solid #6c5ce7; margin: 20px 0; border-radius: 4px; }
        table { width: 100%%; border-collapse: collapse; margin-top: 20px; }
        th { text-align: left; background: #f8f9fa; padding: 12px; font-size: 12px; text-transform: uppercase; color: #636e72; border-bottom: 2px solid #eee; }
        td { padding: 12px; border-bottom: 1px solid #f0f0f0; font-size: 14px; }
        
        /* Mobile Responsiveness */
        @media only screen and (max-width: 600px) {
            .content { padding: 20px; }
            .header { padding: 20px; }
            .header h1 { font-size: 20px; }
            td, th { padding: 8px; font-size: 13px; }
        }
    </style>
</head>
<body>
    <div class="wrapper">
        <br>
        <div class="container">
            <div class="header">
                <h1>BlockVotes</h1>
            </div>
            <div class="content">
                %s
            </div>
            <div class="footer">
                <p>&copy; 2026 BlockVotes Secure Elections.<br>
                This is an automated message. Please do not reply directly.</p>
            </div>
        </div>
        <br>
    </div>
</body>
</html>
`, subject, content)
}

// GenerateOTPEmail creates a beautiful OTP email
func GenerateOTPEmail(otp string) string {
	content := fmt.Sprintf(`
		<h2 style="color: #2d3436; margin-top: 0;">Verify Your Identity</h2>
		<p>Hello,</p>
		<p>You have requested to verify your identity for the BlockVotes platform. Please use the One-Time Password (OTP) below to complete your verification.</p>
		
		<div style="text-align: center; margin: 30px 0;">
			<div style="display: inline-block; background: #dfe6e9; padding: 15px 30px; font-size: 32px; font-weight: bold; color: #2d3436; letter-spacing: 5px; border-radius: 8px; border: 1px dashed #b2bec3;">
				%s
			</div>
		</div>

		<div class="info-box">
			<strong>Note:</strong> This OTP satisfies secure authentication requirements and is valid for <strong>10 minutes</strong>. Do not share this code with anyone.
		</div>

		<p>If you did not request this verification, please ignore this email.</p>
	`, otp)

	return BaseEmailLayout("Your Verification Code", content)
}

// GenerateWelcomeEmail creates a credential email for new voters
func GenerateWelcomeEmail(name, email, password, electionName string) string {
	content := fmt.Sprintf(`
		<h2 style="color: #2d3436; margin-top: 0;">Welcome, %s!</h2>
		<p>Your voter account has been successfully created for the election: <strong>%s</strong>.</p>
		<p>You can now log in to the secure voting portal using the credentials below:</p>

		<div class="info-box" style="background: #fff3cd; border-left-color: #ffc107; color: #856404;">
			<p style="margin: 5px 0;"><strong>Email:</strong> %s</p>
			<p style="margin: 5px 0;"><strong>Temporary Password:</strong> %s</p>
		</div>

		<p>For your security, we strongly recommend that you change your password immediately after your first login.</p>
		
		<div style="text-align: center;">
			<a href="#" class="btn">Login to Vote</a>
		</div>
	`, name, electionName, email, password)

	return BaseEmailLayout("Welcome to BlockVotes", content)
}

// GenerateResultsEmailHTML creates a rich HTML email with results and audit logs
func GenerateResultsEmailHTML(electionName, winnerName string, candidates []map[string]interface{}, logs []AuditLog) string {
	// 1. Build Candidates Table
	var candRows string
	winnerVotes := 0

	for _, c := range candidates {
		name, _ := c["name"].(string)
		votes := 0
		if v, ok := c["voteCount"]; ok {
			switch val := v.(type) {
			case int:
				votes = val
			case int32, int64, float64:
				votes = int(val.(int64)) // simplification, assumes safe cast for display
			}
		}

		if name == winnerName {
			winnerVotes = votes
		}

		row := fmt.Sprintf(`
			<tr>
				<td><strong>%s</strong></td>
				<td>%d</td>
			</tr>
		`, name, votes)
		candRows += row
	}

	// 2. Build Audit Log Table
	var logRows string
	if len(logs) == 0 {
		logRows = `<tr><td colspan="3" style="color:#888; text-align:center;">No activity recorded.</td></tr>`
	} else {
		for _, l := range logs {
			timeStr := l.Timestamp.Format("2006-01-02 15:04:05")
			row := fmt.Sprintf(`
				<tr>
					<td style="color: #666; width: 140px;">%s</td>
					<td><strong>%s</strong></td>
					<td style="color: #555;">%s</td>
				</tr>
			`, timeStr, l.Action, l.Details)
			logRows += row
		}
	}

	content := fmt.Sprintf(`
		<h2 style="color: #2d3436; margin-top: 0;">Election Results Announced</h2>
		<p>The results for <strong>%s</strong> have been finalized. The transparent outcome is presented below.</p>

		<div style="background: #e0f7fa; border-left: 5px solid #00bcd4; padding: 20px; margin: 25px 0; border-radius: 4px;">
			<div style="font-size: 12px; text-transform: uppercase; color: #006064; font-weight: bold; letter-spacing: 1px;">Winner Declared</div>
			<div style="font-size: 24px; color: #006064; margin-top: 5px; font-weight: bold;">%s</div>
			<div style="font-size: 14px; color: #00838f;">with %d verified votes</div>
		</div>

		<h3 style="border-bottom: 2px solid #eee; padding-bottom: 10px; margin-top: 30px;">Vote Count Summary</h3>
		<table>
			<thead><tr><th>Candidate</th><th>Votes</th></tr></thead>
			<tbody>%s</tbody>
		</table>

		<h3 style="border-bottom: 2px solid #eee; padding-bottom: 10px; margin-top: 40px;">Audit Log</h3>
		<p style="font-size: 13px; color: #999;">Immutable record of election events.</p>
		<table>
			<thead><tr><th>Time (UTC)</th><th>Action</th><th>Details</th></tr></thead>
			<tbody>%s</tbody>
		</table>
	`, electionName, winnerName, winnerVotes, candRows, logRows)

	return BaseEmailLayout(fmt.Sprintf("Results: %s", electionName), content)
}
