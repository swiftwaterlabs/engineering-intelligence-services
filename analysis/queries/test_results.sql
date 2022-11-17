select distinct tr.id as Id,
                tr.project as Project,
                element_at(split(tr.project,'.'),1) as Organization,
                element_at(split(tr.project,'.'),2) as Repository,
                from_iso8601_timestamp(tr.analyzedat) as AnalyzedAt,
                cast(tr.metrics.coverage as decimal(10,4)) as CodeCoverage,
                cast(tr.metrics.lines_to_cover as int) as TotalLines,
                cast(tr.metrics.uncovered_lines as int) as UnCoveredLines,
                ro.parentowner as ParentOwner
from engineering_intelligence_prd.test_result tr
         left join engineering_intelligence_prd.repository_owner ro on element_at(split(tr.project,'.'),1) = ro.repository.organization.name AND
                                                                       element_at(split(tr.project,'.'),2) = ro.repository.name
order by tr.project,from_iso8601_timestamp(tr.analyzedat)