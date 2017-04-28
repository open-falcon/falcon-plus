// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
%{

package falcon

import (
	"os"
	"fmt"
)

%}

%union {
	num int
	text string
	b bool
}

%type <b> bool
%type <text> text
%type <num> num

%token <num> NUM
%token <text> TEXT IPA

%token '{' '}' ';'
%token ON YES OFF NO INCLUDE ROOT PID_FILE LOG HOST DISABLED DEBUG
%token CTRL AGENT LOADBALANCE BACKEND
%token UPSTREAM METRIC MIGRATE

%%

config: 
| config conf
;

bool:
  ON  { $$ = true }
| YES { $$ = true }
| OFF { $$ = false }
| NO  { $$ = false }
|     { $$ = true }
;

text:
  IPA  { $$ = string(yy.t) }
| TEXT { $$ = exprText(yy.t) }
;

num:
NUM { $$ = yy.i }
;

ss:
| ss text text ';'    { yy_ss[$2] = $3 }
| ss text num ';'     { yy_ss[$2] = fmt.Sprintf("%d", $3) }
| ss text bool ';'    { yy_ss[$2] = fmt.Sprintf("%v", $3) }
| ss INCLUDE text ';' { yy.include($3) }
;

as:
| as text             { yy_as = append(yy_as, $2) }
| as INCLUDE text ';' { yy.include($3) }
;

conf: ';'
 | PID_FILE text ';' { conf.pidFile = $2 }
 | LOG text num ';'  {
 conf.log = $2
 conf.logv = $3
}| INCLUDE text ';'  { yy.include($2) }
 | ROOT text ';'     { 
 	if err := os.Chdir($2); err != nil {
 		yy.Error(err.Error())
 	}
}| ctrl_mod '}'      {
 	yy_ctrl.Ctrl.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_ctrl.Name = fmt.Sprintf("ctrl_%s", yy_ctrl.Name)
	if yy_ctrl.Host == ""{
		yy_ctrl.Host, _ = os.Hostname()
	}

	if !yy_ctrl.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_ctrl)
	}
}| agent_mod '}'      {
 	yy_agent.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_agent.Name = fmt.Sprintf("agent_%s", yy_agent.Name)
	if yy_agent.Host == ""{
		yy_agent.Host, _ = os.Hostname()
	}

	if !yy_agent.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_agent)
	}
}| loadbalance_mod '}'   {
 	yy_loadbalance.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_loadbalance.Name = fmt.Sprintf("loadbalance_%s", yy_loadbalance.Name)
	if yy_loadbalance.Host == ""{
		yy_loadbalance.Host, _ = os.Hostname()
	}
	if !yy_loadbalance.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_loadbalance)
	}
}| backend_mod '}'      {
 	yy_backend.Configer.Set(APP_CONF_FILE, yy_ss2)
	yy_ss2 = make(map[string]string)

	yy_backend.Name = fmt.Sprintf("backend_%s", yy_backend.Name)
	if yy_backend.Host == ""{
		yy_backend.Host, _ = os.Hostname()
	}
	if !yy_backend.Disabled || yy.debug {
		conf.conf = append(conf.conf, yy_backend)
	}
}
;

////////////////////// ctrl /////////////////////////
ctrl_mod: CTRL text '{' {
	yy_ctrl      = &ConfCtrl{}
	yy_ctrl.Ctrl.Set(APP_CONF_DEFAULT, ConfDefault["ctrl"])
	yy_ctrl.Name = $2
}| ctrl_mod ctrl_mod_item ';'
;

ctrl_mod_item:
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DISABLED bool   { yy_ctrl.Disabled = $2 }
 | DEBUG           { yy_ctrl.Debug = 1 }
 | DEBUG num       { yy_ctrl.Debug = $2 }
 | HOST text       { yy_ctrl.Host = $2 }
 | METRIC '{' as '}' {
 	yy_ctrl.Metrics = yy_as
	yy_as = make([]string, 0)
}| text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }
;

////////////////////// agent /////////////////////////
agent_mod: AGENT text '{' {
	yy_agent      = &ConfAgent{}
	yy_agent.Configer.Set(APP_CONF_DEFAULT, ConfDefault["agent"])
	yy_agent.Name = $2
}| agent_mod agent_mod_item ';'
;

agent_mod_item:
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DISABLED bool   { yy_agent.Disabled = $2 }
 | DEBUG           { yy_agent.Debug = 1 }
 | DEBUG num       { yy_agent.Debug = $2 }
 | HOST text       { yy_agent.Host = $2 }
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | UPSTREAM text   { yy_ss2["upstream"] = $2 }
 | INCLUDE text    { yy.include($2) }
;


////////////////////// loadbalance  /////////////////////////
loadbalance_mod: LOADBALANCE text '{' {
	yy_loadbalance      = &ConfLoadbalance{}
	yy_loadbalance.Configer.Set(APP_CONF_DEFAULT, ConfDefault["loadbalance"])
	yy_loadbalance.Name = $2
}| loadbalance_mod loadbalance_mod_item ';'
;

loadbalance_mod_item:
 | DISABLED bool   { yy_loadbalance.Disabled = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG           { yy_loadbalance.Debug = 1 }
 | DEBUG num       { yy_loadbalance.Debug = $2 }
 | HOST text       { yy_loadbalance.Host = $2 }
 | BACKEND '{' loadbalance_backend '}'
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }

;

loadbalance_backend:
| loadbalance_backend loadbalance_backend_item ';'
;

loadbalance_backend_item:
| text text '{' loadbalance_backend_obj '}' { 
	yy_loadbalance_backend.Type = $1
	yy_loadbalance_backend.Name = $2
	if !yy_loadbalance_backend.Disabled || yy.debug {
		yy_loadbalance.Backend = append(yy_loadbalance.Backend, *yy_loadbalance_backend)
	}
	yy_loadbalance_backend = &LbBackend{}
}
;

loadbalance_backend_obj:
| loadbalance_backend_obj loadbalance_backend_obj_item ';'
;

loadbalance_backend_obj_item:
| DISABLED bool { yy_loadbalance_backend.Disabled = $2 }
| UPSTREAM '{' ss '}' { 
	yy_loadbalance_backend.Upstream = yy_ss
	yy_ss = make(map[string]string)
}
;

////////////////////// backend  /////////////////////////
backend_mod: BACKEND text '{'   {
	yy_backend      = &ConfBackend{}
	yy_backend.Configer.Set(APP_CONF_DEFAULT, ConfDefault["backend"])
	yy_backend.Name = $2
}| backend_mod backend_mod_item ';'
;


backend_mod_item:
 | DISABLED bool   { yy_backend.Disabled = $2 }
 | ROOT text { 
	if err := os.Chdir($2); err != nil {
		yy.Error(err.Error())
	}
}| DEBUG           { yy_backend.Debug = 1 }
 | DEBUG num       { yy_backend.Debug = $2 }
 | HOST text       { yy_backend.Host = $2 }
 | MIGRATE '{' backend_migrate '}'
 | text text       { yy_ss2[$1] = $2 }
 | text num        { yy_ss2[$1] = fmt.Sprintf("%d", $2) }
 | text bool       { yy_ss2[$1] = fmt.Sprintf("%v", $2) }
 | INCLUDE text    { yy.include($2) }
;

backend_migrate:
| backend_migrate backend_migrate_item ';'

backend_migrate_item:
| DISABLED bool { yy_backend.Migrate.Disabled = $2 }
| UPSTREAM '{' ss '}' {
	yy_backend.Migrate.Upstream = yy_ss
	yy_ss = make(map[string]string)
}
;

%%

