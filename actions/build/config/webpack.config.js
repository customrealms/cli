// const path = require('path');

module.exports = {
    module: {
        rules: [
            {
                test: /\.(?:js|ts)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'babel-loader',
                        options: {
                            targets: {
                                browsers: ["ie 7"]
                            },
                            presets: [
                                "@babel/preset-env",
                            ],
                        },
                    },
                    'ts-loader',
                ],
            },
        ],
    },
    resolve: {
        extensions: ['.ts', '.js'],
    },
    output: {
        filename: 'bundle.js',
        environment: {
            arrowFunction: false,
        },
    },
};