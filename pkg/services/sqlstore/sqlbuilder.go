package sqlstore

import (
	"bytes"
	"strings"

	m "github.com/grafana/grafana/pkg/models"
)

type SqlBuilder struct {
	sql    bytes.Buffer
	params []interface{}
}

func (sb *SqlBuilder) Write(sql string, params ...interface{}) {
	sb.sql.WriteString(sql)

	if len(params) > 0 {
		sb.params = append(sb.params, params...)
	}
}

func (sb *SqlBuilder) GetSqlString() string {
	return sb.sql.String()
}

func (sb *SqlBuilder) AddParams(params ...interface{}) {
	sb.params = append(sb.params, params...)
}

func (sb *SqlBuilder) writeDashboardPermissionFilter(user *m.SignedInUser, permission m.PermissionType) {

	if user.OrgRole == m.ROLE_ADMIN {
		return
	}

	okRoles := []interface{}{user.OrgRole}

	if user.OrgRole == m.ROLE_EDITOR {
		okRoles = append(okRoles, m.ROLE_VIEWER)
	}

	falseStr := dialect.BooleanStr(false)

	sb.sql.WriteString(` AND
	(
		dashboard.id IN (
			SELECT distinct Id AS DashboardId FROM (
				SELECT a.Id, coalesce(MAX(case when a.user_id not null then permission end), MAX(case when a.team_id is not null then permission end), MAX(case when a.role is not null then permission end)) as permission
				FROM (
					SELECT d.Id, 
						da.user_id, 
						da.team_id,
						da.role,
						fa.permission as folder_permission, 
						da.permission as dashboard_permission,
						coalesce(da.permission, fa.permission) as permission
					FROM dashboard d
					LEFT JOIN dashboard folder ON folder.Id = d.folder_id
					LEFT JOIN dashboard_acl da ON d.Id = da.dashboard_id OR
					(
						-- include default permissions -->
						da.org_id = -1 AND (
						  (folder.id IS NULL AND d.has_acl = ` + falseStr + `)
						)
					)
					LEFT JOIN dashboard_acl fa ON d.folder_id = fa.dashboard_id OR
					(
						-- include default permissions -->
						da.org_id = -1 AND (
						  (folder.id IS NOT NULL AND folder.has_acl = ` + falseStr + `)
						)
					)
					LEFT JOIN team_member as ugm on ugm.team_id =  da.team_id
					LEFT JOIN org_user ou ON ou.role = da.role AND ou.user_id = ?
					LEFT JOIN org_user ouRole ON ouRole.user_id = ? AND ouRole.org_id = ?
					WHERE
					d.org_id = ? AND
					(
						da.user_id = ? OR
						fa.user_id = ? OR
						ugm.user_id = ? OR
						da.role IN (?` + strings.Repeat(",?", len(okRoles)-1) + `) OR
						fa.role IN (?` + strings.Repeat(",?", len(okRoles)-1) + `)
					)
				) a
				GROUP BY 1
				)
				WHERE permission >= ?
		)
	)`)

	sb.params = append(sb.params, user.UserId, user.UserId, user.OrgId, user.OrgId, user.UserId, user.UserId, user.UserId)
	sb.params = append(sb.params, okRoles...)
	sb.params = append(sb.params, okRoles...)
	sb.params = append(sb.params, permission)
}
