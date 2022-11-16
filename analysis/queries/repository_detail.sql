select r.organization.host as Host,
    r.organization.hosttype as HostType,
    r.organization.name as OrganizationName,
    r.name as Name,
    concat(r.organization.host,'.',r.organization.name,'.',r.name) as FullName,
    r.url as Url,
    r.defaultbranch as DefaultBranch,
    r.createdat as CreatedAt,
    ro.pattern as OwnershipPattern,
    ro.owner as CodeOwner,
    ro.parentowner as ParentOwner,
    case when ro.owner is not null and ro.owner != '' then true else false end HasCodeOwner,
    case when ro.pattern is not null and ro.pattern != '' then true else false end HasCodeOwnerPattern,
    case when ro.parentowner is not null and ro.parentowner != '' then true else false end HasParentOwner
from engineering_intelligence_prd.repository r
left join engineering_intelligence_prd.repository_owner ro on r.id = ro.repository.id
order by r.organization.host,r.organization.name,r.name