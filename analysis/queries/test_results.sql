select id as Id,
       project as Project,
       split(project,'.')[1] as Organization,
       split(project,'.')[2] as Repository,
       from_iso8601_timestamp(analyzedat) as AnalyzedAt,
       cast(metrics.coverage as decimal(10,4)) as CodeCoverage
from engineering_intelligence_prd.test_result