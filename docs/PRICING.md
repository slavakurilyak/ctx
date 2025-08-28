# ctx Pro Pricing (Coming Soon)

## The Insane Math: $1,125 → $10/agent/month

With ctx Pro's 95% token reduction (coming soon), your TOTAL AI costs will become just the $10/agent/month subscription.
That's a 99.1% cost reduction when you factor in the eliminated token costs!

### Cost Comparison
| | Without ctx | With ctx Pro |
|---|---|---|
| **Token Costs** | $1,125/month | ~$0/month |
| **Subscription** | $0 | $10/agent/month |
| **Total Monthly Cost** | $1,125/month | $10/agent/month |
| **Monthly Savings** | - | **$1,115/month** |
| **Annual Savings** | - | **$13,380/year** |

**ctx Pro doesn't just pay for itself - it saves you $1,115/month!**

## Plans & Pricing

**Note: ctx Pro is coming soon! The pricing information below shows what will be available once launched.**

ctx Pro uses simple per-agent pricing with volume discounts for teams:

### Per-Agent Pricing
**$10/agent/month** *(Your entire AI cost after 95% token reduction!)*

Perfect for individual developers and teams of any size who want to enhance their command-line experience with intelligent features.

All features included:
- Pre-tool-use command analysis and optimization
- Post-tool-use insights and recommendations
- Command blocking for dangerous operations
- Token-aware command execution
- Secure credential storage
- All webhook integrations
- Centralized team management
- Shared security policies
- Team-wide command analytics
- Custom webhook configurations
- Priority support
- Audit logs

### Volume Discounts
- **1-4 agents**: $10/agent/month
- **5-9 agents**: $9/agent/month (10% discount)
- **10-24 agents**: $8/agent/month (20% discount)
- **25+ agents**: Contact for custom pricing

## Billing

- **Billing Cycle**: Monthly or Annual (20% discount on annual plans)
- **Payment Methods**: Credit card, PayPal, or invoice (for teams > 10 agents)
- **Currency**: USD (other currencies available upon request)

## Getting Started (Coming Soon)

Once ctx Pro launches, you'll be able to:

1. **Sign up** at [ctx.click](https://ctx.click) (coming soon)
2. **Get your API key** from your dashboard (coming soon)
3. **Login** using `ctx login` with your API key (coming soon)
4. **Check status** with `ctx account` to view your plan details (coming soon)

## Usage Examples

### Single Developer
```bash
# Login with API key
$ echo "your-api-key" | ctx login

# Check account status
$ ctx account
ctx Pro Account Status
======================
Email:       developer@example.com
Tier:        pro
Status:      Active
Valid Until: 2025-02-13

Billing Information:
  Plan:         per-agent
  Price:        $10.00 USD/agent/month
  Agents:       1
  Total:        $10.00 USD/monthly
  Next Billing: 2025-02-13
```

### Team (5 Developers)
```bash
# Login with team API key
$ echo "team-api-key" | ctx login

# Check account status with volume discount
$ ctx account
ctx Pro Account Status
======================
Email:       admin@company.com
Tier:        pro
Status:      Active
Valid Until: 2025-02-13

Billing Information:
  Plan:         per-agent
  Price:        $9.00 USD/agent/month (10% volume discount)
  Agents:       5
  Total:        $45.00 USD/monthly
  Next Billing: 2025-02-13
```

## Webhook Features

Pro accounts unlock powerful webhook integrations:

### Pre-Tool-Use Webhook
- **Command Analysis**: Analyze commands before execution
- **Security Blocking**: Block dangerous commands automatically
- **Command Optimization**: Suggest or apply optimizations
- **Custom Policies**: Apply team-specific rules

### Post-Tool-Use Webhook
- **Output Analysis**: Analyze command results
- **Error Detection**: Identify and explain errors
- **Performance Insights**: Track command performance
- **Learning Recommendations**: Suggest better approaches

## ROI Calculator

### Example: Database-Heavy Workflow
- **Queries per day**: 100
- **Average tokens per raw query**: 2,500
- **Total tokens per month**: 7,500,000

**Without ctx:**
- Cost: ~$1,125/month (at $0.15 per 1k tokens)

**With ctx Pro:**
- Reduced tokens: ~30,000/month (95% reduction)
- Token cost: ~$4.50/month
- ctx Pro subscription: $10/agent/month
- **Total: $14.50/month**
- **Savings: $1,110.50/month (98.7% reduction)**

The subscription pays for itself with just **ONE** optimized database query!

## FAQ

**Q: How does ctx achieve 95% token reduction?**
A: By enabling "measure-then-act" workflows where agents check token counts first, then refine queries to be specific and efficient.

**Q: Can I switch between plans?**
A: Yes, you can upgrade or downgrade at any time. Changes take effect at the next billing cycle.

**Q: Is there a free trial?**
A: Yes, new users get a 14-day free trial of ctx Pro features.

**Q: What happens if my subscription expires?**
A: ctx continues to work with basic features. Pro features (webhooks) are disabled until you renew.

**Q: Can I get a discount for annual billing?**
A: Yes, annual plans receive a 20% discount.

**Q: Do you offer educational discounts?**
A: Yes, we offer 50% discounts for students and educators with valid .edu email addresses.

## Support

For billing questions or support:
- Email: hello@ctx.click
- Documentation: https://docs.ctx.click
- GitHub Issues: https://github.com/slavakurilyak/ctx/issues