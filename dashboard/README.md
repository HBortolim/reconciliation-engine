# Reconciliation Engine Dashboard

A React + TypeScript web dashboard for the reconciliation-engine project, providing comprehensive UI for transaction reconciliation, exception management, fee analysis, and acquirer contract management.

## Project Structure

```
dashboard/
├── src/
│   ├── api/
│   │   ├── client.ts          # Axios instance with auth interceptors
│   │   └── reconciliation.ts  # API endpoints for reconciliation operations
│   ├── components/
│   │   ├── Layout.tsx         # Main layout wrapper with sidebar
│   │   ├── Sidebar.tsx        # Navigation sidebar
│   │   ├── StatusBadge.tsx    # Color-coded status indicators
│   │   ├── MoneyDisplay.tsx   # BRL currency formatter (centavos → R$)
│   │   └── ConfidenceBar.tsx  # Visual confidence score bar (0-100%)
│   ├── pages/
│   │   ├── Dashboard.tsx      # Main dashboard with summary metrics
│   │   ├── ReconciliationRuns.tsx # List of reconciliation runs
│   │   ├── RunDetail.tsx      # Detail view for a single run
│   │   ├── Exceptions.tsx     # Exception management and resolution
│   │   ├── FeeAnalysis.tsx    # Fee variance analysis by acquirer/bandeira
│   │   ├── AcquirerContracts.tsx # Contract CRUD operations
│   │   └── AgingDashboard.tsx # Aging report for unreconciled items
│   ├── types/
│   │   └── index.ts           # Domain types, interfaces, and enums
│   ├── App.tsx                # React Router configuration
│   ├── main.tsx               # React entry point
│   └── index.css              # Tailwind CSS imports
├── index.html                 # HTML entry point
├── vite.config.ts             # Vite build configuration
├── tsconfig.json              # TypeScript configuration
├── tailwind.config.js         # Tailwind CSS configuration
├── postcss.config.js          # PostCSS configuration
├── package.json               # Dependencies and scripts
├── .env.example               # Environment variables template
└── README.md                  # This file
```

## Domain Types

### Enums
- **SourceType**: ACQUIRER, ISSUER, GATEWAY
- **MatchType**: EXACT, FUZZY, MANUAL
- **ExceptionType**: AMOUNT_MISMATCH, MISSING_TRANSACTION, DUPLICATE, DATE_MISMATCH, FEE_DISCREPANCY
- **Severity**: LOW, MEDIUM, HIGH, CRITICAL
- **RunStatus**: PENDING, RUNNING, COMPLETED, FAILED
- **ResolutionStatus**: UNRESOLVED, PENDING_REVIEW, RESOLVED, ESCALATED

### Core Interfaces
- **TransactionRecord**: Individual transaction from a source system
- **ReconciliationPair**: Two matched transactions with confidence score
- **ReconciliationException**: Unmatched or discrepant transaction
- **ReconciliationRun**: Reconciliation execution with statistics
- **AcquirerContract**: Contract terms for an acquirer/bandeira combination
- **FeeSchedule**: Fee structure for a contract

## Features

### Dashboard
- Aggregate metrics: total processed, matched, exceptions, total amount
- Recent reconciliation runs summary table
- Quick status indicators and trend data

### Reconciliation Runs
- List all runs with status, date, match counts
- Clickable rows to view run details
- Trigger new reconciliation runs

### Run Detail
- Statistics and match rate percentage
- Table of matched pairs with confidence scores
- Exceptions list with severity indicators
- Visual confidence bars for match quality

### Exception Management
- Filterable exception table (by severity, status)
- Exception detail modal with resolution workflow
- Add resolution notes and change status
- Track assignment and resolution history

### Fee Analysis
- Period selector (YYYY-MM format)
- Summary cards: contracted vs actual vs variance
- Detailed breakdown by acquirer and bandeira
- Variance highlighting

### Acquirer Contracts
- CRUD operations for contracts
- Display contract effective dates and expiry
- Status indicators (ACTIVE, INACTIVE, PENDING)
- Fee structure associations

### Aging Dashboard
- Unreconciled items bucketed by age (0-7d, 8-15d, 16-30d, 30+d)
- Summary cards with counts and amounts
- Percentage of total calculations
- Detailed aging breakdown table

## Styling

Uses **Tailwind CSS v3.3** for utility-first styling. All components use basic Tailwind classes for responsive design and color-coded status indicators:
- Green: Success/Resolved
- Yellow: Pending/Medium priority
- Orange: High priority/Alert
- Red: Failed/Critical

## API Integration

Configured via environment variables:
- `VITE_API_URL`: Base API URL (default: `http://localhost:8080/api`)

### API Client Features
- Axios instance with configurable base URL
- Auth token injection in request headers
- Automatic 401 handling with redirect to login
- Request/response interceptors ready for extension

### Available Endpoints
- `GET /runs` - List reconciliation runs
- `GET /runs/{id}` - Get run details
- `POST /runs` - Trigger new reconciliation
- `GET /runs/{id}/pairs` - Get matched pairs for a run
- `GET /runs/{id}/exceptions` - Get exceptions for a run
- `PATCH /exceptions/{id}` - Resolve exception
- `GET /contracts` - List acquirer contracts
- `GET /analysis/fees` - Get fee analysis for period

## Installation & Setup

### Prerequisites
- Node.js 18+ (npm 10+)
- React 18.2+
- TypeScript 5.3+

### Steps

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Create environment file**:
   ```bash
   cp .env.example .env
   # Edit .env with your API URL
   ```

3. **Development server**:
   ```bash
   npm run dev
   ```
   Opens on `http://localhost:3000`

4. **Build for production**:
   ```bash
   npm run build
   ```
   Output in `dist/`

5. **Preview production build**:
   ```bash
   npm run preview
   ```

## Technology Stack

- **React 18.2**: UI framework
- **TypeScript 5.3**: Type safety
- **React Router 6.20**: Client-side routing
- **Axios 1.6**: HTTP client with interceptors
- **Tailwind CSS 3.3**: Utility-first styling
- **Recharts 2.10**: Charting library (for future enhancements)
- **Vite 5.0**: Modern build tool
- **PostCSS**: CSS processing with Autoprefixer

## Code Organization

### Separation of Concerns
- **api/**: HTTP client and endpoint definitions
- **types/**: TypeScript domain models and enums
- **pages/**: Full-page components tied to routes
- **components/**: Reusable UI components (Layout, StatusBadge, MoneyDisplay, ConfidenceBar, Sidebar)

### Styling Convention
- Utility-first Tailwind classes
- Component-level scoping via className strings
- Responsive breakpoints: `sm:`, `md:`, `lg:`
- Dark text on light backgrounds for accessibility

### Type Safety
- All API responses typed with domain interfaces
- Props interfaces for all components
- Strict TypeScript mode enabled
- No `any` types without justification

## Future Enhancements

- Charts using Recharts (fee trends, match rate over time)
- Export functionality (CSV, PDF reports)
- Real-time WebSocket updates for run status
- Advanced filtering and search
- User preferences and theme switching
- Detailed audit logs
- Role-based access control (RBAC)
- Dark mode support
- Mobile app version

## Contributing

Follow the established patterns:
1. Create domain types in `src/types/index.ts`
2. Add API functions in `src/api/reconciliation.ts`
3. Create pages in `src/pages/` for routes
4. Create reusable components in `src/components/`
5. Use existing components (StatusBadge, MoneyDisplay, ConfidenceBar)

## License

Same as parent reconciliation-engine project
