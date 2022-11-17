select distinct tr.id as Id,
                tr.host as Host,
                tr.hosttype as HostType,
                tr.project as Project,
                element_at(split(tr.project,'.'),1) as Organization,
                element_at(split(tr.project,'.'),2) as Repository,
                cast(from_iso8601_timestamp(tr.analyzedat) as timestamp) as AnalyzedAt,
                cast(tr.metrics.coverage as decimal(10,4)) as CodeCoverage,
                case when cast(tr.metrics.coverage as decimal(10,4)) between 0 and 10 then 0
                     when cast(tr.metrics.coverage as decimal(10,4)) between 10 and 20 then 10
                     when cast(tr.metrics.coverage as decimal(10,4)) between 20 and 30 then 20
                     when cast(tr.metrics.coverage as decimal(10,4)) between 30 and 40 then 30
                     when cast(tr.metrics.coverage as decimal(10,4)) between 40 and 50 then 40
                     when cast(tr.metrics.coverage as decimal(10,4)) between 50 and 60 then 50
                     when cast(tr.metrics.coverage as decimal(10,4)) between 60 and 70 then 60
                     when cast(tr.metrics.coverage as decimal(10,4)) between 70 and 80 then 70
                     when cast(tr.metrics.coverage as decimal(10,4)) between 80 and 90 then 80
                     when cast(tr.metrics.coverage as decimal(10,4)) between 90 and 100 then 90
                     when cast(tr.metrics.coverage as decimal(10,4)) >=1000 then 100
                     else null
                    end CodeCoverageRange,
                cast(tr.metrics.lines_to_cover as int) as TotalLines,
                cast(tr.metrics.uncovered_lines as int) as UnCoveredLines,
                ro.parentowner as ParentOwner
from engineering_intelligence_prd.test_result tr
         left join engineering_intelligence_prd.repository_owner ro on lower(element_at(split(tr.project,'.'),1)) = lower(ro.repository.organization.name) AND
                                                                       lower(element_at(split(tr.project,'.'),2)) = lower(ro.repository.name)
order by tr.project