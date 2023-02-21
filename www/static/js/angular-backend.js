
'use strict';

var BackendService;

angular.module( 'backend', [] )

.factory( '$backend', [ '$http', '$timeout', '$rootScope', function( $http, $timeout, $rootScope ) {
    return function( script ) {
        return new BackendService( $http, $timeout, $rootScope, script );
    };
} ] );

// ////////////////////////////////////////////////// 
// BackendService

BackendService = function( http, timeout, scope, script ) {
    this._http = http;
    this._timeout = timeout;
    this._scope = scope;
    this._script = script;
    this._success = null;
    this._error = null;
};

BackendService.prototype = {

    // ////////////////////////////////////////////////// 
    // action
    
    action: function( action, parameters ) {
        var self = this;
        
        // build url 
        
        var url = 'http://backend.amnezic.com/' + self._script,
            sep = '?';
        parameters = parameters || {};
        parameters[ 'output' ] = 'jsonp';
        parameters[ 'callback' ] = 'JSON_CALLBACK';
        parameters[ 'action' ] = action;
        if ( self._scope.auth ) {
            parameters[ 'uid' ] = self._scope.auth.uid;
            parameters[ 'session' ] = self._scope.auth.session;
        }
        parameters[ '_' ] = Date.now();
        for ( var key in parameters ) {
            if ( !parameters.hasOwnProperty( key ) ) {
                continue;
            }
            if ( !parameters[ key ] ) {
                continue;
            }
            url += sep + key + '=' + ( parameters[ key ] || '' );
            sep = '&';
        }
        
        console.log( ' >>> ' + url );
        
        // execute url
        
        self._scope.$broadcast( 'error' );
        self._scope.$broadcast( 'loading', true );
        self._http.jsonp( url, { cache: false } )
        .success( function( data ) {
            self._scope.$broadcast( 'loading', false );
            if ( data.success ) {

                // success
                
                if ( self._success ) {
                    delete data.success;
                    self._success( data );
                }
            }
            else {
            
                // error
                
                if ( self._error ) {
                    delete data.success;
                    self._error( data );
                } else {
                    self._scope.$broadcast( 'error', data.error ? data.error : 'internal.error' );
                }
            }
        } )
        .error( function() {
            self._scope.$broadcast( 'loading', false );
            if ( self._error ) {
                self._error();
            } else {
                self._scope.$broadcast( 'error', 'internal.error' );
            }
        } );
        return self;
    },
    
    // ////////////////////////////////////////////////// 
    // callbacks
    
    success: function( callback ) {
        this._success = callback;
        return this;
    },
    
    error: function( callback ) {
        this._error = callback;
        return this;
    }
  
};        
