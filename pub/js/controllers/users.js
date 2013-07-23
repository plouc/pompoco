function UsersCtrl($scope, $http) {
  $scope.creating = false;
  $scope.editing  = false;

  $http.get('http://localhost:8080/users')
  .then(function(response) {
    $scope.users = response.data;
  });

  $scope.create = function() {
    $scope.cancel();
    $scope.creating = true;
  };

  $scope.edit = function(user) {
    $scope.creating = false;
    $scope.editing  = true;
    $scope.user     = user;
  };

  $scope.cancel = function() {
    $scope.creating = false;
    $scope.editing  = false;
    $scope.user = {
      'name':          '',
      'github_active': false,
      'github_name':   '',
      'gitlab_active': false,
      'gitlab_name':   '',
    };
  };

  $scope.remove = function(user) {
    $http.delete('http://localhost:8080/user/' + user._id)
    .then(function(response) {
      $scope.users = response.data;
    });
  };

  $scope.save = function() {
    var method = 'POST',
        url    = 'http://localhost:8080/user';
    if ($scope.editing === true) {
      method = 'PUT';
      url    += '/' + $scope.user._id;
    }

    $http({
      'method':  method,
      'url':     url,
      'data':    serialize($scope.user),
      'headers': {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })
    .then(function(response) {
      $scope.users = response.data;
      $scope.cancel();
    });
  };
}
  