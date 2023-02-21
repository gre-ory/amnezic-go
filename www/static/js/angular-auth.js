
'use strict';

var AuthService;

angular.module( 'auth', [ 'storage', 'backend' ] )

.factory( '$auth', [ '$backend', function( $backend ) {
    return new AuthService( $backend( 'security/access.py' ) );
} ] );

// ////////////////////////////////////////////////// 
// AuthService

AuthService = function( backend ) {
    this._backend = backend;
    this._scope = backend._scope;
};

AuthService.prototype = {
        
    // ////////////////////////////////////////////////// 
    // login
    
    login: function( email, password ) {
        var self = this;
        
        self._backend.action( 'login', { email: email, password: password } )
        .success( function( data ) {
            console.log( 'login user ' + email + ' with id ' + data.uid + ' and session ' + data.session );
            self._scope.$broadcast( 'auth', { email: email, connected: true, uid: data.uid, session: data.session } );
        } );
    },
    
    // ////////////////////////////////////////////////// 
    // logout
    
    logout: function( email ) {
        var self = this;
        
        self._backend.action( 'logout', { email: email } )
        .success( function( data ) {
            console.log( 'logout user ' + email );
            self._scope.$broadcast( 'auth', { email: email, connected: false } );
        } );
    }
    
};

// ////////////////////////////////////////////////// 
// helper

var auth_on = [ '$storage', '$rootScope', '$auth', function( $storage, $rootScope, $auth ) {
    console.log( 'auth_on' );
    console.log( $rootScope.auth );
    if ( !$rootScope.auth ) {
        $rootScope.auth = { email: $storage.get( 'email' ), connected: false };
    }
    console.log( $rootScope.auth );
    $rootScope.login = function() {
        console.log( 'login' );
        console.log( $rootScope.auth );
        $auth.login( $rootScope.auth.email, $rootScope.auth.password );
    };
    $rootScope.logout = function() {
        console.log( 'logout' );
        console.log( $rootScope.auth );
        $auth.logout( $rootScope.auth.email );
    };
    $rootScope.$on( 'auth', function( event, auth ) {
        console.log( 'on.auth' );
        console.log( auth );
        console.log( $rootScope.auth );
        $rootScope.auth = auth;
        console.log( $rootScope.auth );
        console.log( $rootScope.auth.email );
        $storage.add( 'email', $rootScope.auth.email );
        console.log( $storage.get( 'email' ) );
        $rootScope.$broadcast( 'refresh' );
    } );
} ];        
