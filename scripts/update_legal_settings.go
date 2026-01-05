package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

var privacyPolicy = `# Privacy Policy

**Effective Date:** January 6, 2026
**Last Updated:** January 6, 2026

## Introduction

Yi-API ("we," "our," or "us") respects your privacy and is committed to protecting your personal information. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our AI API gateway service at api.ai80.vip (the "Service").

## Information We Collect

### Information You Provide
- **Account Information:** Email address, username, and password when you register
- **Payment Information:** Billing details processed securely through our payment processors (Stripe, Creem, WeChat Pay, Alipay)
- **API Usage Data:** API keys, request logs, and usage statistics
- **Communication Data:** Messages you send to our customer support

### Automatically Collected Information
- **Technical Data:** IP address, browser type, device information, operating system
- **Usage Data:** API call logs, response times, error rates, feature usage
- **Cookies:** Session cookies for authentication and functionality

## How We Use Your Information

We use your information to:
- Provide, maintain, and improve the Service
- Process your API requests and manage your account
- Handle billing and payment processing
- Monitor usage and prevent abuse
- Respond to customer support inquiries
- Send important service updates and notifications
- Comply with legal obligations

## Data Sharing and Disclosure

We do NOT sell your personal information. We may share data with:

- **Payment Processors:** Stripe, Creem, WeChat Pay, Alipay for payment processing
- **AI Service Providers:** Your API requests are forwarded to third-party AI providers (OpenAI, Anthropic, Google, etc.) as necessary to fulfill your requests
- **Legal Requirements:** When required by law or to protect our rights
- **Business Transfers:** In connection with mergers, acquisitions, or asset sales

## Data Security

We implement industry-standard security measures including:
- Encryption of data in transit (TLS/HTTPS)
- Secure password hashing
- Regular security audits
- Access controls and authentication
- Secure API key management

However, no system is 100% secure. You are responsible for keeping your account credentials confidential.

## Data Retention

We retain your data for as long as:
- Your account is active
- Needed to provide the Service
- Required by law or for legitimate business purposes
- You can request account deletion at any time

API request logs are typically retained for 90 days for analytics and debugging purposes.

## Your Rights

Depending on your location, you may have the right to:
- Access your personal information
- Correct inaccurate data
- Delete your account and data
- Export your data
- Opt-out of marketing communications
- Object to certain data processing

To exercise these rights, contact us at: yian20133213@gmail.com

## Third-Party Services

Our Service integrates with multiple third-party AI providers. Each provider has their own privacy policy:
- OpenAI: https://openai.com/privacy
- Anthropic (Claude): https://www.anthropic.com/privacy
- Google (Gemini): https://policies.google.com/privacy
- And others

When you use our Service, your data may be processed by these providers according to their policies.

## International Data Transfers

Your information may be transferred to and processed in countries other than your own. We ensure appropriate safeguards are in place for such transfers.

## Children's Privacy

Our Service is not intended for users under 13 years of age. We do not knowingly collect information from children.

## Changes to This Policy

We may update this Privacy Policy periodically. We will notify you of material changes by:
- Posting the updated policy on our website
- Updating the "Last Updated" date
- Sending email notifications for significant changes

## Contact Us

If you have questions about this Privacy Policy, please contact us:

**Email:** yian20133213@gmail.com
**Website:** https://api.ai80.vip

---

© 2026 Yi-API (QuantumNous). All rights reserved.`

var termsOfService = `# Terms of Service

**Effective Date:** January 6, 2026
**Last Updated:** January 6, 2026

## Acceptance of Terms

By accessing or using Yi-API ("Service") at api.ai80.vip, you agree to be bound by these Terms of Service ("Terms"). If you do not agree to these Terms, do not use the Service.

## Description of Service

Yi-API is an AI API gateway and relay service that provides unified access to multiple AI service providers including OpenAI, Anthropic, Google, DeepSeek, and others. We offer:
- API request routing and load balancing
- Quota and billing management
- Multi-channel integration
- Usage analytics and monitoring

## Account Registration

### Eligibility
- You must be at least 13 years old to use the Service
- You must provide accurate and complete registration information
- You are responsible for maintaining the confidentiality of your account credentials
- You may not share your account with others

### Account Security
- You are responsible for all activities under your account
- Notify us immediately of any unauthorized access
- We reserve the right to suspend accounts showing suspicious activity

## API Usage and Restrictions

### Acceptable Use
You agree to use the Service only for lawful purposes and in accordance with these Terms.

### Prohibited Activities
You may NOT:
- Use the Service for illegal activities
- Attempt to circumvent rate limits or quotas
- Reverse engineer or attempt to extract source code
- Resell access without explicit authorization
- Generate spam, phishing, or malicious content
- Violate intellectual property rights
- Abuse, harass, or harm others
- Overload or disrupt the Service infrastructure
- Use the Service to train competing AI models
- Share API keys publicly or with unauthorized parties

### API Rate Limits
- Rate limits apply based on your account tier
- Excessive usage may result in throttling or suspension
- Enterprise plans are available for higher limits

## Billing and Payment

### Pricing
- Pricing is based on token usage and model selection
- Current pricing is available on our website
- We reserve the right to modify pricing with notice

### Payment Terms
- Prepaid quota system requires advance payment
- Payments are non-refundable except as required by law
- We accept payments via Stripe, Creem, WeChat Pay, and Alipay
- Failed payments may result in service suspension

### Refunds
- Unused quota may be refunded at our discretion
- Refund requests must be submitted within 30 days
- Refunds may take 5-10 business days to process

## Intellectual Property

### Your Content
- You retain ownership of content you submit
- You grant us a license to process and transmit your content to fulfill requests
- You represent that you have rights to all content you submit

### Our Service
- The Service, including software, trademarks, and content, is owned by Yi-API
- We grant you a limited, non-exclusive license to use the Service
- You may not copy, modify, or create derivative works

## Data and Privacy

Your use of the Service is also governed by our Privacy Policy, which is incorporated into these Terms by reference. Key points:
- We collect and process data as described in our Privacy Policy
- Your API requests are forwarded to third-party AI providers
- We implement security measures but cannot guarantee absolute security
- You are responsible for complying with data protection laws applicable to your use

## Service Availability

### Uptime
- We strive to maintain high availability but do not guarantee 100% uptime
- Scheduled maintenance will be announced in advance when possible

### Service Changes
- We may modify or discontinue features with or without notice
- We will provide reasonable notice for significant changes

## Termination

### By You
- You may terminate your account at any time
- Contact us to request account deletion

### By Us
We may suspend or terminate your account if:
- You violate these Terms
- You engage in fraudulent activity
- You fail to pay required fees
- Required by law
- We discontinue the Service

Upon termination:
- Your access will be immediately revoked
- Unused quota may be forfeited
- We may delete your data per our retention policy

## Disclaimers and Limitation of Liability

### Disclaimer of Warranties
THE SERVICE IS PROVIDED "AS IS" WITHOUT WARRANTIES OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, OR NON-INFRINGEMENT.

### Limitation of Liability
TO THE MAXIMUM EXTENT PERMITTED BY LAW, YI-API SHALL NOT BE LIABLE FOR:
- Indirect, incidental, special, or consequential damages
- Loss of profits, data, or business opportunities
- Service interruptions or data loss
- Third-party AI provider errors or outages

OUR TOTAL LIABILITY SHALL NOT EXCEED THE AMOUNT YOU PAID IN THE PAST 12 MONTHS.

## Indemnification

You agree to indemnify and hold Yi-API harmless from claims, damages, and expenses (including legal fees) arising from:
- Your use of the Service
- Your violation of these Terms
- Your violation of any third-party rights
- Content you submit through the Service

## Third-Party Services

The Service relies on third-party AI providers. We are not responsible for:
- Third-party service availability or performance
- Third-party terms of service or privacy policies
- Content generated by third-party AI models

## Dispute Resolution

### Governing Law
These Terms are governed by the laws of [Your Jurisdiction], without regard to conflict of law provisions.

### Dispute Process
1. Contact us to resolve disputes informally: yian20133213@gmail.com
2. If unresolved, disputes may be subject to binding arbitration
3. You may have the right to opt-out of arbitration in certain jurisdictions

## General Provisions

### Entire Agreement
These Terms, together with our Privacy Policy, constitute the entire agreement between you and Yi-API.

### Severability
If any provision is found invalid, the remaining provisions remain in effect.

### No Waiver
Our failure to enforce any right does not waive that right.

### Assignment
We may assign these Terms; you may not without our consent.

### Updates to Terms
We may modify these Terms at any time. Continued use after changes constitutes acceptance.

## Contact Information

For questions about these Terms:

**Customer Support Email:** yian20133213@gmail.com
**Website:** https://api.ai80.vip

---

© 2026 Yi-API (QuantumNous). All rights reserved.`

func main() {
	common.SetupGinLog()
	common.SysLog("正在初始化数据库...")
	err := model.InitDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			log.Fatalf("关闭数据库失败: %v", err)
		}
	}()

	// Initialize configuration
	config.GlobalConfig.Init()

	// Load current legal settings
	legalSettings := system_setting.GetLegalSettings()

	// Update with new content
	legalSettings.PrivacyPolicy = privacyPolicy
	legalSettings.UserAgreement = termsOfService

	// Convert to JSON
	jsonData, err := json.Marshal(legalSettings)
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}

	// Update database
	err = model.UpdateOption("legal", string(jsonData))
	if err != nil {
		log.Fatalf("更新数据库配置失败: %v", err)
	}

	fmt.Println("✅ 隐私政策和服务条款已成功更新到数据库")
	fmt.Println("📝 隐私政策长度:", len(privacyPolicy), "字符")
	fmt.Println("📝 服务条款长度:", len(termsOfService), "字符")
	fmt.Println("")
	fmt.Println("🔗 访问以下链接查看:")
	fmt.Println("   隐私政策: https://api.ai80.vip/privacy")
	fmt.Println("   服务条款: https://api.ai80.vip/user-agreement")
}
