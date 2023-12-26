const HtmlWebpackPlugin = require('html-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const path = require('path');

module.exports = {
    entry: {
        Todo: './src/Todo.ts',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
            //   {
            //     test:/\.html$/,
            //     use:{
            //         loader:'file-loader',
            //         options:{
            //             name:'[name].[ext]'
            //         }
            //     },
            //     exclude: path.resolve(__dirname,'src/index.html')
            //   }
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'static/js'),
    },
    devServer: {
        static: path.join(__dirname, 'static/js'),
        devMiddleware: {
            writeToDisk: true,
        },
        compress: true,
        port: 4000,
    },
    plugins: [
        new Dotenv(),
    ]
};
