/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const { exit, listenerCount } = require('process');
const yargs = require('yargs');
const { showHidden } = require('yargs');

const ccpPath = path.resolve(__dirname, '..', 'network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const argv = yargs
        .command(
            'upload',
            'upload file.', 
            yargs => {
                yargs
                    .option('file', {
                        alias: 'f',
                        description: 'File Path',
                        type: 'string',
                    })
                    .option('key', {
                        alias: 'k',
                        description: 'The key which upload to the ledger',
                        type: 'string',
                    })    
                    .option('tag', {
                        alias: 't',
                        description: 'Tag',
                        type: 'string',
                        default: 'default'
                    })
                    .demandOption(['file', 'key'])
            }
        )
        .command(
            'download',
            'download file.', 
            yargs => {
                yargs
                    .option('file', {
                        alias: 'f',
                        description: 'File Path',
                        type: 'string',
                    })
                    .option('key', {
                        alias: 'k',
                        description: 'The key which upload to the ledger',
                        type: 'string',
                    })    
                    .demandOption(['file', 'key'])
            }
        )
        .command(
            'list',
            'get file list.', 
            yargs => {
                yargs
                    .option('tag', {
                        alias: 't',
                        description: 'Tag',
                        type: 'string',
                    })
            }
        )
        .command(
            'show',
            'show file.', 
            yargs => {
                yargs
                    .option('key', {
                        alias: 'k',
                        description: 'The key which upload to the ledger',
                        type: 'string',
                    })    
                    .demandOption(['key'])
            }
        )
        .demandCommand(1, 'You need at least one command')
        .help()
        .alias('help', 'h')
        .argv;

async function main() {   
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        //console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('vote');

        if (argv._ == 'upload') {
            await upload(contract)
        } else if (argv._ == 'download') {
            await download(contract)
        } else if (argv._ == 'list') {
            await list(contract)
        } else if (argv._ == 'show') {
            await show(contract)
        }

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Error occured: ${error}`);
        process.exit(1);
    }
}

async function upload(contract) {
    console.log('read file...');
    var data = fs.readFileSync(argv.file)
    var buff = new Buffer(data);
    var encodedData = buff.toString('base64');

    console.log('upload file...');
    try {
        const result = await contract.submitTransaction('upload', argv.key, argv.tag, encodedData);
        console.log(`file uploaded as "${argv.key}"`);
    } catch (error) {
        console.log(`cannot upload file. check if key:"${argv.key}" already exist`);
    }
} 

async function download(contract) {
    console.log('download file...');
    try {
        const result = await contract.evaluateTransaction('download', argv.key);

        console.log('save file...');
        var encodedData = result.toString();
        var buff = Buffer.from(encodedData, 'base64');
        fs.writeFileSync(argv.file, buff);
            
        console.log(`file saved as "${argv.file}"`);
    }
    catch (error) {
        console.log(`cannot download file. check if key:"${argv.key}" exist`);
    }
}

async function list(contract) {
    console.log('getting file list...');
    const result = await contract.evaluateTransaction('list');
    var obj = JSON.parse(result);
    var index = 1
        
    console.log('-------------------------------------------------------------------------------'); 
    for (var i=0; i<obj.length; i++) {
        if (!argv.tag || obj[i].Tag == argv.tag) {
            console.log(`${index++}. key: ${obj[i].Key}, tag: ${obj[i].Tag}`);
        }
    }
    if (index == 1) console.log(`no such files with tag:${argv.tag}`)
    console.log('-------------------------------------------------------------------------------');
}

async function show(contract) {
    console.log('getting file...');

    try {
        const result = await contract.evaluateTransaction('show', argv.key);
        var obj = JSON.parse(result);
        console.log(`key: ${obj.Key}, tag: ${obj.Tag}`)
    } catch (error) {
        console.log(`cannot get file. check if key:"${argv.key}" exist`);
    }
}

main();
