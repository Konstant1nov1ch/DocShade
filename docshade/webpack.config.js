const webpack = require('webpack');
const path = require('path');

module.exports = {
  entry: './src/index.tsx', // или ваш путь к входному файлу
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      },
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env.REACT_APP_BACKEND_HOST': JSON.stringify(process.env.REACT_APP_BACKEND_HOST)
    })
  ],
};
