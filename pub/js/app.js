var pompoco = angular.module('pompoco', ['ui.state']);

serialize = function(obj, prefix) {
    var str = [];
    for(var p in obj) {
        var k = prefix ? prefix + "[" + p + "]" : p, v = obj[p];
        str.push(typeof v == "object" ? 
            serialize(v, k) :
            encodeURIComponent(k) + "=" + encodeURIComponent(v));
    }
    return str.join("&");
}



pompoco.config(function($stateProvider, $routeProvider) {
    $stateProvider
    .state('index', {
        url: '', // root route
        controller: 'GitlabCtrl',
        views: {
            "subnav": {
                templateUrl: "views/subnav.html"
            },
            "content": {
                templateUrl: "views/content.html"
            }
        }
    })

    .state('settings', {
      'url':         '/settings',
      'templateUrl': 'views/settings.html',
      'controller':  function($scope, $http) {
        $http.get('http://localhost:8080/settings')
        .then(function(response) {
          console.log(response.data);
        });
      }
    })

    // Users
    .state('users', {
      'url':         '/users',
      'templateUrl': 'views/users/users.html',
      'controller':  'UsersCtrl'
    })
    .state('users.create', {
      'url':         '/new',
      'templateUrl': 'views/users/form.html',
      'controller':  function($scope, $http) {
        $scope.title = 'New user';
        $scope.save = function() {
          $http({
            'method':  'POST',
            'url':     'http://localhost:8080/user',
            'data':    serialize($scope.user),
            'headers': {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          })
          .then(function(response) {
            $scope.users = response.data;
          });
        };
      }
    })
    .state('users.user', {
      'url':         '/:id',
      'templateUrl': 'views/users/user.html',
      'controller':  function($scope, $http, $stateParams) {

        $scope.refreshing = false;

        $http.get('http://localhost:8080/user/' + $stateParams.id)
        .then(function(response) {
          $scope.user = response.data;
        });

        $scope.consolidate = function() {
          $scope.refreshing = true;
          $http.get('http://localhost:8080/user/' + $stateParams.id + '/consolidate')
          .then(function(response) {
            $scope.user = response.data;
            $scope.refreshing = false;
          });
        };
      }
    })
    .state('users.edit', {
      'url':         '/:id/edit',
      'templateUrl': 'views/users/form.html',
      'controller':  function($scope, $http, $state, $stateParams) {

        $scope.title  = '';
        $scope.loaded = false;

        $http.get('http://localhost:8080/user/' + $stateParams.id)
        .then(function(response) {
          $scope.user = response.data;
          $scope.title = 'Edit user ' + $scope.user.name;
          $scope.loaded = true;
        });

        $scope.save = function() {
          if ($scope.loaded === false) {
            return;
          }
          $http({
            'method':  'PUT',
            'url':     'http://localhost:8080/user/' + $scope.user._id,
            'data':    serialize($scope.user),
            'headers': {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          })
          .then(function(response) {
            $scope.projects = response.data;
            $state.transitionTo('users.user', { 'id': $scope.user._id });
          });
        };
      }
    })
    .state('users.delete', {
      'url':         '/:id/delete',
      'templateUrl': 'views/users/delete.html',
      'controller':  function($scope, $http, $state, $stateParams) {
        $http.get('http://localhost:8080/user/' + $stateParams.id)
        .then(function(response) {
          $scope.user = response.data;
        });

        $scope.delete = function() {
          $http.delete('http://localhost:8080/user/' + $scope.user._id)
          .then(function(response) {
            $scope.projects = response.data;
            $state.transitionTo('users');
          });
        };
      }
    })

    // Projects
    .state('projects', {
      'url':         '/projects',
      'templateUrl': 'views/projects/projects.html',
      'controller':  function($scope, $http) {
        $http.get('http://localhost:8080/projects')
        .then(function(response) {
          $scope.projects = response.data;
        });
      }
    })
    .state('projects.create', {
      'url':         '/new',
      'templateUrl': 'views/projects/form.html',
      'controller':  function($scope, $http) {
        $scope.title = 'New project';
        $scope.project = {
          'name':        '',
          'description': ''
        };
        $scope.save = function() {
          $http({
            'method':  'POST',
            'url':     'http://localhost:8080/project',
            'data':    serialize($scope.project),
            'headers': {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          })
          .then(function(response) {
            $scope.projects = response.data;
            $scope.cancel();
          });
        };
      }
    })
    .state('projects.project', {
      'url':         '/:id',
      'templateUrl': 'views/projects/project.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/project/' + $stateParams.id)
        .then(function(response) {
          $scope.project = response.data;
        });
      }
    })
    .state('projects.edit', {
      'url':         '/:id/edit',
      'templateUrl': 'views/projects/form.html',
      'controller':  function($scope, $http, $state, $stateParams) {

        $scope.title  = '';
        $scope.loaded = false;

        $http.get('http://localhost:8080/project/' + $stateParams.id)
        .then(function(response) {
          $scope.project = response.data;
          $scope.title = 'Edit project ' + $scope.project.name;
          $scope.loaded = true;
        });

        $scope.save = function() {
          if ($scope.loaded === false) {
            return;
          }
          $http({
            'method':  'PUT',
            'url':     'http://localhost:8080/project/' + $scope.project._id,
            'data':    serialize($scope.project),
            'headers': {
              'Content-Type': 'application/x-www-form-urlencoded'
            }
          })
          .then(function(response) {
            $scope.projects = response.data;
            $state.transitionTo('projects');
          });
        };
      }
    })
    .state('projects.delete', {
      'url':         '/:id/delete',
      'templateUrl': 'views/projects/delete.html',
      'controller':  function($scope, $http, $state, $stateParams) {
        $http.get('http://localhost:8080/project/' + $stateParams.id)
        .then(function(response) {
          $scope.project = response.data;
        });

        $scope.delete = function() {
          $http.delete('http://localhost:8080/project/' + $scope.project._id)
          .then(function(response) {
            $scope.projects = response.data;
            $state.transitionTo('projects');
          });
        };
      }
    })

    .state('timeline', {
      'url':         '/timeline',
      'templateUrl': 'views/timeline.html',
      'controller':  function($scope, $http) {

        $scope.synchronizing = false;

        $http.get('http://localhost:8080/events')
        .then(function(response) {
          $scope.events = response.data;
        });

        $scope.syncEvents = function() {
          $scope.synchronizing = true;
          $http.get('http://localhost:8080/events/sync')
          .then(function(response) {
            $scope.events = response.data;
            $scope.synchronizing = false;
          });
        };
      }
    })

    // Github
    .state('github', {
      'url':         '/github',
      'templateUrl': 'views/github/github.html',
      'controller':  function($scope, $http) {
        $http.get('http://localhost:8080/github/repos')
        .then(function(response) {
          $scope.repos = response.data;
        });
      }
    })
    .state('github.usage', {
      'url':         '/usage',
      'templateUrl': 'views/github/usage.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/github/usage')
        .then(function(response) {
          console.log(response.data);
          $scope.usage = response.data;
        });
      }
    })
    .state('github.repo', {
      'url':         '/repo/:owner/:name',
      'templateUrl': 'views/github/repo.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/github/repo/' + $stateParams.owner + '/' + $stateParams.name)
        .then(function(response) {
          $scope.repo = response.data;
        });
      }
    })

    // Gitlab
    .state('gitlab', {
      'url':         '/gitlab',
      'templateUrl': 'views/gitlab/gitlab.html',
      'controller':  function($scope, $http) {
        $http.get('http://localhost:8080/gitlab/projects')
        .then(function(response) {
          $scope.projects = response.data;
        });
      }
    })
    .state('gitlab-timeline', {
      'url':         '/gitlab/timeline',
      'templateUrl': 'views/gitlab/timeline.html',
      'controller':  function($scope, $http) {
        $http.get('http://localhost:8080/gitlab/timeline')
        .then(function(response) {
          $scope.activity = response.data;
        });
      }
    })
    .state('gitlab.project', {
      'url':         '/project/:id',
      'templateUrl': 'views/gitlab/project.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/gitlab/project/' + $stateParams.id)
        .then(function(response) {
          $scope.project = response.data;
        });
      }
    })
    .state('gitlab.project.branches', {
      'url':         '/branches',
      'templateUrl': 'views/gitlab/branches.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/gitlab/project/' + $stateParams.id + '/branches')
        .then(function(response) {
          $scope.branches = response.data;
        });
      }
    })
    .state('gitlab.project.tags', {
      'url':         '/tags',
      'templateUrl': 'views/gitlab/tags.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/gitlab/project/' + $stateParams.id + '/tags')
        .then(function(response) {
          $scope.tags = response.data;
        });
      }
    })


    // Jira
    .state('jira', {
      'url':         '/jira/timeline',
      'templateUrl': 'views/jira/timeline.html',
      'controller':  function($scope, $http, $stateParams) {
        $http.get('http://localhost:8080/jira/timeline')
        .then(function(response) {
          $scope.activity = response.data;
        })
      }
    })
})
.run(function ($rootScope, $state, $stateParams) {
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
});



function MainCtrl($scope, $http) {

  $scope.currentUser    = null;
  $scope.bindingProject = false;
  $scope.projectBinding = {
    'type':   null,
    'object': {}
  };

  $http.get('http://localhost:8080/users')
  .then(function(response) {
    $scope.users = response.data;
  });

  $http.get('http://localhost:8080/projects')
  .then(function(response) {
    $scope.projects = response.data;
  });

  $scope.bindProject = function(type, object) {
    console.log(type, object);
    $scope.bindingProject = true;
    $scope.projectBinding.type   = type;
    $scope.projectBinding.object = object;
    $scope.projectBinding.role   = '';
    switch (type) {
      case 'user':
        $scope.bindTitle = 'attach user ' + object.name + ' to project';
      case 'gitlab-project':
        $scope.bindTitle = 'attach gitlab project ' + object.name + ' to project';
    }
  };
  $scope.submitProjectBinding = function() {
    console.log($scope.projectBinding);
  };
  $scope.cancelProjectBinding = function() {
    $scope.bindingProject = false;
  }
}