import Keycloak from 'keycloak-js';

// Setup Keycloak instance as needed
// Pass initialization options as required or leave blank to load from 'keycloak.json'
// const keycloak = new Keycloak()

// http://localhost:9999/realms/myrealm/
// http://localhost:9999/realms/myrealm/protocol/openid-connect/login-status-iframe.html/init?client_id=myclient&origin=http%3A%2F%2Flocalhost%3A3000
const initOptions = {
    url: 'http://localhost:9999',
    realm: 'myrealm',
    clientId: 'myclient',
    redirectUri: 'http://localhost:3000//login/auth/keycloak',
    KeycloakResponseType: 'code',
    "public-client": true,
    "confidential-port": 0
}

// export default keycloak = new Keycloak(initOptions);

var keycloak = new Keycloak(initOptions);


// keycloak.init({ onLoad: initOptions.onLoad, KeycloakResponseType: 'code' }).success((auth) => {
//     if (!auth) {
//         window.location.reload();
//     } else {
//         console.info("Authenticated");
//     }
//     setTimeout(() => {
//         keycloak.updateToken(70).success((refreshed) => {
//             if (refreshed) {
//                 console.debug('Token refreshed' + refreshed);
//             } else {
//                 console.warn('Token not refreshed, valid for '
//                     + Math.round(keycloak.tokenParsed.exp + keycloak.timeSkew - new Date().getTime() / 1000) + ' seconds');
//             }
//         }).error(() => {
//             console.error('Failed to refresh token');
//         });

//     }, 60000)
// }).error(() => {
//     console.error("Authenticated Failed");
// });



export default keycloak
