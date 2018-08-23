class OutlyerCli < Formula
  desc "Outlyer CLI allows to easily manage your Outlyer account via command line."
  homepage "https://www.outlyer.com/"
  url "https://github.com/outlyerapp/outlyer-cli/releases/download/0.2.1/outlyer-cli_0.2.1_Darwin_x86_64.tar.gz"
  version "0.2.1"
  sha256 "3bdeaa77a8d62c95dbe440086a6a2704e2018710ccd7b7b72f8f09526dadd690"

  def install
    bin.install "outlyer"
  end
end
