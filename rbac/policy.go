package rbac

import "easyca/model"

//添加权限
func AddPolicy(user, url, del string) (bool, error) {
	return globalEnforcer.AddPolicy(user, url, del)
}

//移除权限
func RemovePolicy(user, url, del string) (bool, error) {
	return globalEnforcer.RemovePolicy(user, url, del)
}

//获取所有角色
func GetAllSubjects() []string {
	//rLock.RLock()
	//defer rLock.RUnlock()
	return globalEnforcer.GetAllSubjects()
}

//获取角色的权限
func GetPermissionsForUser(user string) [][]string {
	//rLock.RLock()
	//defer rLock.RUnlock()
	return globalEnforcer.GetPermissionsForUser(user)
}

//编辑角色的权限
func ModifyPermissionsForUser(param *model.ModifyRole) (err error) {
	//rLock.Lock()
	//defer rLock.Unlock()
	//删除所有用户角色
	_, err = globalEnforcer.DeletePermissionsForUser(param.User)
	if err != nil {
		return err
	}
	//添加用户角色
	for _, role := range param.Roles {
		_, err = AddPolicy(param.User, role, param.Del)
		if err != nil {
			return err
		}
	}
	return
}