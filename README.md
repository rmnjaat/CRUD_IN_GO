# Go-SaaS Foundation

A production-ready, scalable architecture for multi-tenant SaaS applications built with Go and PostgreSQL. This project provides foundational modules for structured logging and intelligent database provisioning, designed to support enterprise-grade SaaS applications with strong tenant isolation and operational observability.

## 🏗️ Project Overview

This repository implements a distributed multi-tenant architecture with two core modules:

1. **Custom Logger Module** - A fully custom-built logging library with structured output for application monitoring across microservices
2. **Multi-Tenant Database Provisioning & Migration System** - An intelligent database management tool that dynamically provisions isolated PostgreSQL databases per tenant using a catalog database pattern

The architecture prioritizes:
- 🔒 **Strong data isolation** between tenants
- 📊 **Structured observability** for monitoring and debugging
- ⚡ **Scalable tenant onboarding** with automated provisioning
- 🔄 **Automated schema migrations** for consistency
- ♻️ **Reusable components** across microservices

---

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Logger Module](#logger-module)
- [Multi-Tenant Database System](#multi-tenant-database-system)
- [Repository Structure](#repository-structure)
- [Getting Started](#getting-started)
- [Future Roadmap](#future-roadmap)
- [Engineering Highlights](#engineering-highlights)

---

## 🏛️ Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Go-SaaS Foundation                          │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────────────────────┐
│  Microservices   │         │      Custom Logger Module        │
│  (e.g., Users)   │────────▶│  - File-based output             │
│                  │         │  - Console output                │
│  - REST APIs     │         │  - Structured logs               │
│  - Business      │         │  - 5 severity levels             │
│    Logic         │         │  - Runtime context capture       │
└──────────────────┘         └──────────────────────────────────┘
        │
        │
        ▼
┌─────────────────────────────────────────────────────────────────┐
│            Multi-Tenant Database Architecture                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────┐         ┌─────────────────────────┐  │
│  │  Catalog Database   │         │   Tenant Databases      │  │
│  │  ═════════════════  │         │   ═══════════════════   │  │
│  │                     │         │                         │  │
│  │  ┌───────────────┐  │         │  ┌──────────────────┐  │  │
│  │  │ Tenants Table │  │────────▶│  │  Tenant A DB     │  │  │
│  │  ├───────────────┤  │         │  │  - users table   │  │  │
│  │  │ • tenant_id   │  │         │  │  - employees     │  │  │
│  │  │ • db_name     │  │         │  │  - versioned     │  │  │
│  │  │ • db_user     │  │         │  └──────────────────┘  │  │
│  │  │ • db_password │  │         │                         │  │
│  │  │ • is_active   │  │         │  ┌──────────────────┐  │  │
│  │  │ • timestamps  │  │         │  │  Tenant B DB     │  │  │
│  │  └───────────────┘  │         │  │  - users table   │  │  │
│  │                     │         │  │  - employees     │  │  │
│  └─────────────────────┘         │  │  - versioned     │  │  │
│            │                     │  └──────────────────┘  │  │
│            │                     │                         │  │
│            ▼                     │  ┌──────────────────┐  │  │
│  ┌─────────────────────┐         │  │  Tenant N DB     │  │  │
│  │  Migration Engine   │─────────▶  │  - users table   │  │  │
│  │  (Goose Framework)  │         │  │  - employees     │  │  │
│  │                     │         │  │  - versioned     │  │  │
│  │  - Catalog DB       │         │  └──────────────────┘  │  │
│  │    migrations       │         │                         │  │
│  │  - Tenant DB        │         └─────────────────────────┘  │
│  │    migrations       │                                      │
│  └─────────────────────┘                                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📝 Logger Module

### Overview

A **fully custom-built** logging library (not a wrapper) designed for structured application monitoring across distributed microservices. The logger captures runtime context and provides multiple severity levels for comprehensive observability.

### Design Goals

- ✅ **Custom Implementation**: Built from scratch without wrapping existing libraries
- ✅ **Structured Logging**: Standardized format for parsing and analysis
- ✅ **Runtime Context**: Automatic capture of file, line number, and timestamp
- ✅ **Multiple Outputs**: File-based persistence and console output
- ✅ **Microservice-Ready**: Reusable across multiple services
- ✅ **Observability Foundation**: Designed for future integration with monitoring platforms

### Severity Levels

The logger supports five severity levels:

| Level   | Method         | Use Case                                    |
|---------|----------------|---------------------------------------------|
| INFO    | `Info(msg)`    | General informational messages              |
| DEBUG   | `Debug(msg)`   | Detailed debugging information              |
| WARNING | `Warning(msg)` | Warning messages for potential issues       |
| ERROR   | `Error(msg)`   | Error conditions that need attention        |
| FATAL   | `Fatal(msg)`   | Critical failures affecting availability    |

### Log Format

Each log entry follows a structured format:

```
[MODULE_NAME] [TIMESTAMP] [LEVEL] [FILENAME:LINE] MESSAGE
```

**Example:**
```
[UserService] [2024-11-25 14:30:45] [INFO] [main.go:43] Server started at port :8080
[UserService] [2024-11-25 14:30:50] [ERROR] [database.go:127] Failed to connect to database
```

### Output Targets

#### 1. File-Based Output
- Logs are written to configurable file paths
- Automatic directory creation
- Append-only mode for log rotation compatibility
- Permissions: `0644` (read/write for owner, read for others)

#### 2. Console Output
- Real-time console output for development and debugging
- Simultaneous file and console logging

#### 3. Database Persistence (Planned)
- Infrastructure exists for database logging
- Future implementation for centralized log aggregation

### Technical Implementation

- **Runtime Context Capture**: Uses `runtime.Caller(2)` to automatically capture calling file and line number
- **UTC Timestamps**: All logs use UTC timezone for consistency across distributed systems
- **Thread-Safe**: Safe for concurrent use across goroutines
- **Zero Dependencies**: No external logging library dependencies

### File Structure

```
packages/logger/
├── main.go       # Core logger struct and severity level methods
├── formatter.go  # Log formatting and runtime context capture
└── writer.go     # File I/O operations and directory management
```

---

## 🗄️ Multi-Tenant Database System

### Overview

An intelligent database migration system built on the **Goose framework** that implements a **catalog database pattern** for dynamic provisioning and management of isolated PostgreSQL databases per tenant. This design ensures strong data separation while maintaining scalability and automated schema consistency.

### System Design

#### Catalog Database Pattern

The system uses a centralized **catalog database** to track tenant metadata and database locations:

The catalog database stores tenant metadata including tenant identifiers, database names for isolation, credentials, and activation status.

### Tenant Lifecycle Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Tenant Onboarding Flow                       │
└─────────────────────────────────────────────────────────────────┘

1. New Tenant Registration
         │
         ▼
   ┌─────────────────┐
   │  Add tenant to  │
   │  catalog DB     │
   │  'tenants' table│
   └────────┬────────┘
            │
            ▼
2. Migration Engine Startup
         │
         ▼
   ┌─────────────────┐
   │  Query catalog  │
   │  DB for all     │
   │  tenant DBs     │
   └────────┬────────┘
            │
            ▼
3. Dynamic Provisioning
         │
         ▼
   ┌─────────────────┐
   │ For each tenant:│
   │  - Check if DB  │
   │    exists       │
   │  - Create if    │
   │    not exists   │
   └────────┬────────┘
            │
            ▼
4. Schema Migration
         │
         ▼
   ┌─────────────────┐
   │  Run Goose      │
   │  migrations on  │
   │  tenant DB:     │
   │  - users table  │
   │  - employees    │
   │  - etc.         │
   └────────┬────────┘
            │
            ▼
5. Ready for Use
   ┌─────────────────┐
   │  Tenant DB is   │
   │  provisioned    │
   │  and versioned  │
   └─────────────────┘
```

### Migration Strategy

#### Two-Tier Migration System

**1. Catalog Database Migrations** (`catalog_db_migrations/`)
- Creates and manages the `tenants` table
- Run once on system initialization
- Manages central metadata schema

**2. Tenant Database Migrations** (`tenant_db_migration/`)
- Applied to each tenant's isolated database
- Standard application schema (users, employees, etc.)
- Version-controlled with Goose
- Ensures schema consistency across all tenants



### Data Isolation Strategy

#### Database-Level Isolation

**Strong Isolation Benefits:**
- ✅ Complete physical data separation
- ✅ Independent database backups and restores
- ✅ Per-tenant performance tuning
- ✅ Compliance-friendly (GDPR, HIPAA)
- ✅ Prevents cross-tenant data leakage

**Why This Design?**

| Approach | Pros | Cons | Our Choice |
|----------|------|------|------------|
| Shared DB, Shared Schema | Simple | ❌ Security risk, complex queries | ❌ Not used |
| Shared DB, Schema per Tenant | Moderate isolation | ⚠️ Connection limits | ⚠️ Alternative |
| **Database per Tenant** | **Strong isolation** | Higher resource usage | ✅ **Selected** |

**Scalability Considerations:**
- PostgreSQL can handle hundreds of databases on modern hardware
- Each tenant database is lightweight (only their data)
- Connection pooling manages resource usage
- Horizontal scaling via database sharding (future enhancement)

### Technical Implementation

**Migration Execution:**
- Goose framework for version control
- SQL-based migrations for transparency
- Automatic rollback support
- Migration history tracked per database

### Configuration

**Environment Variables** (`.env`):
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
CATALOG_DB_NAME=saas_catalog
```

### File Structure

```
migration/
├── main.go                      # Migration orchestrator
├── catalog_db_migrations/       # Catalog DB schema
│   └── 001_tenant_table.sql     # Creates tenants table
└── tenant_db_migration/         # Tenant schema templates
    ├── 001_users.sql            # Users table
    └── 002_employees.sql        # Employees table
```

### Scaling Characteristics

**Current Capacity:**
- Supports hundreds of tenants on single PostgreSQL instance
- Each tenant gets isolated database with full schema
- Automated provisioning eliminates manual setup

**Future Scaling:**
- Database sharding across multiple PostgreSQL servers
- Catalog database clustering for high availability
- Read replicas for tenant databases
- Geographic distribution for global tenants

---

## 📁 Repository Structure

```
Distributed_multi_tenant_architecture/
│
├── packages/
│   └── logger/                  # Custom Logger Module
│       ├── main.go              # Logger struct and severity methods
│       ├── formatter.go         # Log formatting and context capture
│       ├── writer.go            # File I/O operations
│       └── go.mod               # Module dependencies
│
├── migration/                   # Database Migration System
│   ├── main.go                  # Migration orchestrator
│   ├── catalog_db_migrations/   # Catalog database schema
│   │   └── 001_tenant_table.sql
│   ├── tenant_db_migration/     # Tenant database schema templates
│   │   ├── 001_users.sql
│   │   └── 002_employees.sql
│   ├── .env                     # Database configuration
│   ├── go.mod
│   └── go.sum
│
└── Users/                       # Example Microservice
    ├── main.go                  # HTTP server setup
    ├── routes.go                # API routing
    ├── common/
    │   ├── handle_health.go     # Health check endpoint
    │   └── routes.go
    ├── .env
    ├── go.mod
    └── go.sum
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 14 or higher
- Goose CLI (optional, for manual migrations)

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/yourusername/Distributed_multi_tenant_architecture.git
cd Distributed_multi_tenant_architecture
```

2. **Set up environment variables:**

Create `.env` files in `migration/` directory:
```bash
# migration/.env
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
CATALOG_DB_NAME=saas_catalog
```

3. **Install dependencies:**
```bash
cd migration
go mod download

cd ../packages/logger
go mod download

cd ../../Users
go mod download
```



---

## 🔮 Future Roadmap

### Real-Time Observability Platform

Building on the custom logger and database infrastructure, the next phase involves creating a comprehensive observability platform:

#### 1. Custom Logger Service Integration
- Centralized log aggregation across all microservices
- Database persistence for log entries
- Advanced querying and filtering capabilities
- Log correlation across distributed systems

#### 2. Database Query Metrics Collection
- Automatic query performance tracking
- SQL execution time monitoring
- Connection pool metrics
- Transaction latency measurement

#### 3. API Performance Tracking
- Request/response time visualization
- Endpoint performance benchmarking
- Slow endpoint detection and alerting
- Throughput and concurrency metrics

#### 4. Slow Query Detection
- Real-time identification of slow database queries
- Query optimization recommendations
- Historical query performance trends
- Per-tenant query analysis

#### 5. Response Time Visualization
- Interactive dashboards for latency metrics
- Percentile-based analysis (p50, p95, p99)
- Time-series visualization
- Anomaly detection and alerting

#### 6. System Health Monitoring
- Critical system health indicators
- Resource utilization tracking (CPU, memory, disk)
- Service availability monitoring
- Automated health checks and alerts

**Technology Stack (Planned):**
- Time-series database (e.g., TimescaleDB, InfluxDB)
- Visualization layer (e.g., Grafana)
- Alert management (e.g., AlertManager)
- Distributed tracing (e.g., OpenTelemetry)

---

## 💡 Engineering Highlights

### Production-Ready Design

✅ **Scalable Multi-Tenant Architecture**
- Database-level isolation ensures security and compliance
- Catalog pattern enables dynamic tenant onboarding
- Horizontal scaling capabilities for growing tenant base

✅ **Custom-Built Observability**
- Logger built from scratch without third-party dependencies
- Structured logging format for machine parsing
- Foundation for comprehensive monitoring platform

✅ **Automated Database Management**
- Zero-touch tenant provisioning with Goose framework
- Automated schema migrations prevent drift
- Idempotent operations safe for re-execution

✅ **Microservice-Ready Components**
- Reusable logger module across services
- Standardized logging interface
- Modular architecture for independent deployment

✅ **Strong Data Isolation**
- Physical database separation per tenant
- Eliminates cross-tenant data leakage
- Compliance-friendly (GDPR, HIPAA, SOC 2)

✅ **Developer Experience**
- Simple API for logging (`log.Info()`, `log.Error()`, etc.)
- Automatic context capture (file, line, timestamp)
- Clear migration strategy and documentation


