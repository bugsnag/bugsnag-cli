import 'package:flutter/material.dart';

void main() {
  runApp(const App());
}

class App extends StatelessWidget {
  const App({Key? key}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Barebones Template',
      theme: ThemeData(
        primarySwatch: Colors.grey,
      ),
      home: const Page(),
    );
  }
}

/// The Model Widget defines external parameters.
class Page extends StatefulWidget {
  const Page({Key? key}) : super(key: key);

  @override
  PageController createState() => PageController();
}

/// The Controller handles the state and business logic.
class PageController extends State<Page> {
  @override
  Widget build(BuildContext context) => PageView(this);
}

/// The View is purely declarative.
class PageView extends StatelessWidget {
  final PageController state;
  const PageView(this.state, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const Scaffold();
  }
}
